package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	uuid "github.com/satori/go.uuid"
	"github.com/thteam47/go-identity-authen-api/errutil"
	v1 "github.com/thteam47/go-identity-authen-api/pkg/api-client/identity-api"
	"github.com/thteam47/go-identity-authen-api/pkg/db"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	"github.com/thteam47/go-identity-authen-api/pkg/pb"
	"github.com/thteam47/go-identity-authen-api/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ComponentContanier struct {
	userService          UserRepository
	authenInfoRepository AuthenInfoRepository
	jwtRepository        JwtRepository
	authRepository       *grpcauth.AuthInterceptor
	handler              *db.Handler
}

func getUser(item *v1.User) (*models.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &models.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

var errorCodeBadRequest = 400

func NewComponentContanier(userService UserRepository, authenInfoRepository AuthenInfoRepository, jwtRepository JwtRepository, authRepository *grpcauth.AuthInterceptor, handler *db.Handler) *ComponentContanier {
	return &ComponentContanier{
		userService:          userService,
		authenInfoRepository: authenInfoRepository,
		jwtRepository:        jwtRepository,
		authRepository:       authRepository,
		handler:              handler,
	}
}

func (inst *ComponentContanier) Login(userContext grpcauth.UserContext, ctx context.Context, reqCtx *pb.Context, loginType string, username string, password string, otp int32, requestId string, typeMfa string) (string, int, string, error) {
	tokenInfo := &models.TokenInfo{}
	var user *v1.User
	if loginType == "UsernamePassword" {
		userItem, err := inst.userService.FindByLoginName(userContext, username)
		if err != nil {
			return "", 0, "", errutil.Wrapf(err, "userService.FindByLoginName")
		}
		if userItem == nil {
			return "", errorCodeBadRequest, "Username or password incorrect", nil
		}
		if userItem.Status == "pending" {
			return "", errorCodeBadRequest, "Account is not activated", nil
		}
		if userItem.Status == "verified" {
			return "", errorCodeBadRequest, "Account is not approved", nil
		}
		authenInfo, err := inst.authenInfoRepository.GetOneByAttr(userContext, map[string]string{
			"user_id": userItem.UserId,
		})
		if err != nil {
			return "", 0, "", errutil.Wrapf(err, "authenInfoRepository.GetOneByAttr")
		}
		if authenInfo == nil {
			return "", 0, "", nil
		}

		isComparePassword := util.CompareHashPassword(authenInfo.HashPassword, password)

		if !isComparePassword {
			return "", errorCodeBadRequest, "Username or password incorrect", nil
		}
		if len(authenInfo.Mfas) == 0 {
			tokenInfo.AuthenticationDone = true
		} else {
			enabledMfa := false
			for _, item := range authenInfo.Mfas {
				if item.Enabled {
					enabledMfa = true
				}
			}
			if !enabledMfa {
				tokenInfo.AuthenticationDone = true
			}
		}
		user = userItem
	} else if loginType == "AccessToken" {
		userContext, err := inst.authRepository.Authentication(ctx, reqCtx, "@any", "@any")
		if err != nil {
			return "", errorCodeBadRequest, "", status.Errorf(codes.PermissionDenied, errutil.Message(err))
		}
		tokenInfo = userContext.GetTokenInfo()
		requestIdCache := ""
		err = inst.handler.RedisRepository.GetValueCache(fmt.Sprintf("request-id-%s", userContext.GetUserId()), &requestIdCache)
		if err != nil {
			return "", errorCodeBadRequest, "Request Id expired. Please login again", nil
		}
		if requestId == "" {
			return "", errorCodeBadRequest, "Request Id Not Found. Please login again", nil
		}
		if strings.TrimSpace(requestId) != requestIdCache {
			return "", errorCodeBadRequest, "Request Id expired. Please login again", nil
		}
		verifyMfa, err := inst.verifyMfa(userContext, typeMfa, otp)
		if err != nil {
			return "", errorCodeBadRequest, "inst.verifyMfa", nil
		}
		err = inst.handler.RedisRepository.RemoveValueCache(fmt.Sprintf("request-id-%s", userContext.GetUserId()))
		if err != nil {
			return "", errorCodeBadRequest, "Request Id expired. Please login again", nil
		}
		if typeMfa == "EmailOtp" {
			err = inst.handler.RedisRepository.RemoveValueCache(fmt.Sprintf("email-otp-%s", userContext.GetUserId()))
			if err != nil {
				return "", errorCodeBadRequest, "Request Id expired. Please login again", nil
			}
		}
		if !verifyMfa {
			return "", errorCodeBadRequest, "Invalid Otp", nil
		}
		tokenInfo.AuthenticationDone = true
	} else {
		return "", errorCodeBadRequest, "Login type unavailable", nil
	}

	if user != nil {
		tokenInfo.PermissionAll = user.PermissionAll
		tokenInfo.UserId = user.UserId
		for _, key := range user.Permissions {
			permission := &models.Permission{
				Privilege: key.Privilege,
				Actions:   key.Actions,
			}
			tokenInfo.Permissions = append(tokenInfo.Permissions, permission)
		}
		tokenInfo.Role = user.Role
	}
	if tokenInfo.Role == "admin" {
		tokenInfo.Subject = "admin"
	}
	tokenInfo.Exp = int32(time.Now().Add(inst.handler.Exp).Unix())
	token, err := inst.jwtRepository.Generate(tokenInfo)
	if err != nil {
		return "", errorCodeBadRequest, "jwtRepository.Generate", nil
	}
	return token, 0, "", nil
}

func (inst *ComponentContanier) PrepareLogin(userContext grpcauth.UserContext) (string, string, string, string, string, []string, string, error) {
	if userContext.IsAuthenDone() {
		return userContext.GetAccessToken(), "", "", "", "", nil, "", nil
	}
	authenInfo, err := inst.authenInfoRepository.GetOneByAttr(userContext, map[string]string{
		"user_id": userContext.GetUserId(),
	})
	if err != nil {
		return "", "", "", "", "", nil, "", errutil.Wrapf(err, "authenInfoRepository.GetOneByAttr")
	}
	if authenInfo == nil {
		return userContext.GetAccessToken(), "", "", "", "", nil, "", nil
	}

	requestId := uuid.NewV4().String()

	err = inst.handler.RedisRepository.SetValueCache(fmt.Sprintf("request-id-%s", userContext.GetUserId()), requestId, inst.handler.TimeRequestId)
	if err != nil {
		return "", "", "", "", "", nil, "", errutil.Wrapf(err, "RedisRepository.SetValueCache")
	}
	availbableMfas := []string{}
	secret := ""
	url := ""
	message := ""
	typeMfa := ""
	for _, item := range authenInfo.Mfas {
		if item.Type == "Totp" {
			if item.Enabled {
				if !item.Configured {
					secret = item.Secret
					url = item.Url
					message = "Please add your TOTP to your OTP Application now"
				} else {
					message = "Please enter your OTP from your OTP Application"
				}
				typeMfa = "Totp"
				break
			}
		} else if item.Type == "EmailOtp" {
			if item.Enabled {
				otp, err := util.GenerateCodeOtp()
				if err != nil {
					return "", "", "", "", "", nil, "", errutil.Wrapf(err, "generateCodeOtp")
				}
				err = inst.handler.RedisRepository.SetValueCache(fmt.Sprintf("email-otp-%s", userContext.GetUserId()), otp, inst.handler.TimeEmailOtp)
				if err != nil {
					return "", "", "", "", "", nil, "", errutil.Wrapf(err, "RedisRepository.SetValueCache")
				}
				err = util.SendMail([]string{item.PublicData}, fmt.Sprintf("Your OTP code is %d", otp))
				if err != nil {
					return "", "", "", "", "", nil, "", errutil.Wrapf(err, "util.SendMail")
				}
				message = "OTP sent to your email"
				typeMfa = "EmailOtp"
				break
			}
		}
	}
	for _, item := range authenInfo.Mfas {
		if item.Enabled {
			availbableMfas = append(availbableMfas, item.Type)
		}
	}
	return userContext.GetAccessToken(), requestId, secret, url, message, availbableMfas, typeMfa, nil
}

func (inst *ComponentContanier) verifyMfa(userContext grpcauth.UserContext, typeMfa string, otp int32) (bool, error) {
	otpCache := 0
	if typeMfa == "EmailOtp" {
		err := inst.handler.RedisRepository.GetValueCache(fmt.Sprintf("email-otp-%s", userContext.GetUserId()), &otpCache)
		if err != nil {
			return false, nil
		}
		if int32(otpCache) != otp {
			return false, nil
		}
		return true, nil
	}
	if typeMfa == "Totp" {
		authenInfo, err := inst.authenInfoRepository.GetOneByAttr(userContext, map[string]string{
			"user_id": userContext.GetUserId(),
		})
		if err != nil {
			return false, errutil.Wrapf(err, "authenInfoRepository.GetOneByAttr")
		}
		secret := ""
		configured := false
		index := -1
		mfa := &models.Mfa{}
		for key, item := range authenInfo.Mfas {
			if item.Type == "Totp" {
				mfa = item
				secret = item.Secret
				configured = item.Configured
				index = key
			}
		}
		valid := totp.Validate(strconv.Itoa(int(otp)), secret)
		if !configured && valid {
			mfa.Configured = true
			authenInfo.Mfas[index] = mfa
			err = inst.authenInfoRepository.UpdateOneByAttr(userContext, userContext.GetUserId(), map[string]interface{}{
				"mfas": authenInfo.Mfas,
			})
			if err != nil {
				return false, errutil.Wrapf(err, "authenInfoRepository.UpdateOneByAttr")
			}
		}
		return valid, nil
	}
	return false, nil
}

func (inst *ComponentContanier) Logout(userContext grpcauth.UserContext) error {
	err := inst.handler.RedisRepository.SetValueCache(fmt.Sprintf("invalid-token-%s", userContext.GetUserId()), userContext.GetAccessToken(), inst.handler.Exp)
	if err != nil {
		return errutil.Wrapf(err, "RedisRepository.SetValueCache")
	}
	return nil
}
