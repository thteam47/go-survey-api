package grpcapp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thteam47/go-identity-authen-api/errutil"
	"github.com/thteam47/go-identity-authen-api/pkg/db"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	pb "github.com/thteam47/go-identity-authen-api/pkg/pb"
	"github.com/thteam47/go-identity-authen-api/pkg/repository"
	"github.com/thteam47/go-identity-authen-api/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityAuthenService struct {
	pb.IdentityAuthenServiceServer
	handler              *db.Handler
	componentContanier   *repository.ComponentContanier
	userService          repository.UserRepository
	authRepository       *grpcauth.AuthInterceptor
	authenInfoRepository repository.AuthenInfoRepository
}

func NewIdentityAuthenService(handler *db.Handler, componentContanier *repository.ComponentContanier, userService repository.UserRepository, authRepository *grpcauth.AuthInterceptor, authenInfoRepository repository.AuthenInfoRepository) *IdentityAuthenService {
	return &IdentityAuthenService{
		handler:              handler,
		componentContanier:   componentContanier,
		authRepository:       authRepository,
		authenInfoRepository: authenInfoRepository,
		userService:          userService,
	}
}

func getMfa(item *pb.Mfa) (*models.Mfa, error) {
	if item == nil {
		return nil, nil
	}
	user := &models.Mfa{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getMfas(items []*pb.Mfa) ([]*models.Mfa, error) {
	mfas := []*models.Mfa{}
	for _, item := range items {
		mfa, err := getMfa(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getMfas")
		}
		mfas = append(mfas, mfa)
	}
	return mfas, nil
}

func makeMfa(item *models.Mfa) (*pb.Mfa, error) {
	mfa := &pb.Mfa{}
	err := util.ToMessage(item, mfa)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return mfa, nil
}

func makeMfas(items []*models.Mfa) ([]*pb.Mfa, error) {
	mfas := []*pb.Mfa{}
	for _, item := range items {
		mfa, err := makeMfa(item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeMfa")
		}
		mfas = append(mfas, mfa)
	}
	return mfas, nil
}

func (inst *IdentityAuthenService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginName := strings.TrimSpace(strings.ToLower(req.Username))
	loginType := strings.TrimSpace(req.Type)
	password := strings.TrimSpace(req.Password)
	requestId := strings.TrimSpace(req.RequestId)
	typeMfa := strings.TrimSpace(req.TypeMfa)
	if loginType == "UsernamePassword" && (loginName == "" || password == "") {
		return &pb.LoginResponse{
			Token:     "",
			ErrorCode: 400,
			Message:   "Username or password incorrect",
		}, nil
	}
	token, errorCode, message, err := inst.componentContanier.Login(nil, ctx, req.Ctx, loginType, loginName, password, req.Otp, requestId, typeMfa)
	if err != nil {
		return nil, errutil.Wrapf(err, "componentContanier.Login")
	}
	return &pb.LoginResponse{
		Token:     token,
		ErrorCode: int32(errorCode),
		Message:   message,
	}, nil
}

func (inst *IdentityAuthenService) PrepareLogin(ctx context.Context, req *pb.PrepareLoginRequest) (*pb.PrepareLoginResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}

	token, requestId, secret, url, message, availbableMfas, typeMfa, err := inst.componentContanier.PrepareLogin(userContext)
	if err != nil {
		return nil, errutil.Wrapf(err, "componentContanier.PrepareLogin")
	}

	return &pb.PrepareLoginResponse{
		Token:         token,
		Message:       message,
		RequestId:     requestId,
		AvailableMfas: availbableMfas,
		TypeMfa:       typeMfa,
		Secret:        secret,
		Url:           url,
	}, nil
}

func (inst *IdentityAuthenService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.MessageResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-authen-api:authen-info", "update")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	// }

	userId := strings.TrimSpace(req.UserId)
	password := strings.TrimSpace(req.Password)
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errutil.Wrapf(err, "util.HashPassword")
	}
	err = inst.authenInfoRepository.UpdateOneByAttr(nil, userId, map[string]interface{}{
		"hash_password": hashPassword,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.UpdateOneByAttr(")
	}
	return &pb.MessageResponse{
		Ok:      true,
		Message: "Update passord successful",
	}, nil

}

func (inst *IdentityAuthenService) Logout(ctx context.Context, req *pb.Request) (*pb.MessageResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.componentContanier.Logout(userContext)
	if err != nil {
		return nil, errutil.Wrapf(err, "componentContanier.Logout")
	}
	return &pb.MessageResponse{}, nil
}

func (inst *IdentityAuthenService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.MessageResponse, error) {
	loginName := strings.TrimSpace(strings.ToLower(req.Data))
	message, err := inst.authenInfoRepository.ForgotPassword(nil, loginName)
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.ForgotPassword")
	}
	return &pb.MessageResponse{
		Message: message,
		Ok: true,
	}, nil
}

func (inst *IdentityAuthenService) UpdateForgotPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.MessageResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	password := strings.TrimSpace(req.Password)
	_, err = inst.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		Ctx:      req.Ctx,
		UserId:   userContext.GetUserId(),
		Password: password,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "inst.UpdatePassword")
	}
	fmt.Println(userContext.GetAccessToken())

	err = inst.handler.RedisRepository.SetValueCache(fmt.Sprintf("invalid-token-%s", userContext.GetAccessToken()), userContext.GetAccessToken(), 5*time.Minute)
	if err != nil {
		return nil, errutil.Wrapf(err, "RedisRepository.SetValueCache")
	}
	return &pb.MessageResponse{
		Message: "Update Password Successfull",
		Ok: true,
	}, nil
}

func (inst *IdentityAuthenService) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.MessageResponse, error) {
	username := strings.TrimSpace(strings.ToLower(req.Username))
	email := strings.TrimSpace(strings.ToLower(req.Email))
	fullName := strings.TrimSpace(strings.ToLower(req.FullName))
	message, err := inst.authenInfoRepository.RegisterUser(nil, username, fullName, email)
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.RegisterUser")
	}
	return &pb.MessageResponse{
		Message: message,
		Ok:      true,
	}, nil
}

func (inst *IdentityAuthenService) VerifyUser(ctx context.Context, req *pb.Request) (*pb.MessageResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.userService.VerifyUser(userContext, "")
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.RegisterUser")
	}

	return &pb.MessageResponse{
		Message: "Account has been actived",
		Ok:      true,
	}, nil
}
func (inst *IdentityAuthenService) VerifyForgotPassword(ctx context.Context, req *pb.Request) (*pb.MessageResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	if userContext == nil {
		return &pb.MessageResponse{
			Ok: false,
		}, nil
	}
	return &pb.MessageResponse{
		Ok: true,
	}, nil
}
