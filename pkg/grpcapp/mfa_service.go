package grpcapp

import (
	"context"
	"strings"

	"github.com/thteam47/go-identity-authen-api/errutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	pb "github.com/thteam47/go-identity-authen-api/pkg/pb"
	"github.com/thteam47/go-identity-authen-api/util"
	"golang.org/x/exp/slices"
)

func (inst *IdentityAuthenService) GetMfaType(ctx context.Context, req *pb.StringRequest) (*pb.MfaResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-authen-api:authen-info", "update")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	// }

	userId := strings.TrimSpace(req.Value)

	authenInfo, err := inst.authenInfoRepository.GetOneByAttr(nil, map[string]string{
		"user_id": userId,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.GetOneByAttr(")
	}

	mfas, err := makeMfas(authenInfo.Mfas)
	if err != nil {
		return nil, errutil.Wrapf(err, "makeMfas")
	}
	return &pb.MfaResponse{
		Mfas: mfas,
	}, nil

}

func (inst *IdentityAuthenService) UpdateMfa(ctx context.Context, req *pb.UpdateMfaRequest) (*pb.MessageResponse, error) {
	// userContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "identity-authen-api:authen-info", "update")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	// }

	userId := strings.TrimSpace(req.UserId)
	mfas, err := getMfas(req.Mfas)
	if err != nil {
		return nil, errutil.Wrapf(err, "getMfas")
	}
	mfasValid := []*models.Mfa{}
	typeMfas := []string{"Totp", "EmailOtp"}
	for _, item := range mfas {
		if slices.Contains(typeMfas, strings.TrimSpace(item.Type)) {
			if strings.TrimSpace(item.Type) == "Totp" {
				hash, err := util.HashPassword(userId)
				if err != nil {
					return nil, errutil.Wrapf(err, "util.HashPassword")
				}
				key, err := util.GenerateTotp(hash)
				if err != nil {
					return nil, errutil.Wrapf(err, "util.GenerateTotp")
				}
				item.Secret = key.Secret()
				if err != nil {
					return nil, errutil.Wrapf(err, "util.HashPassword")
				}
				item.Url = key.URL()
			}
			item.PublicData = strings.TrimSpace(strings.ToLower(item.PublicData))
			mfasValid = append(mfasValid, item)
		}
	}

	if len(mfasValid) == 0 {
		return &pb.MessageResponse{}, nil
	}

	err = inst.authenInfoRepository.UpdateOneByAttr(nil, userId, map[string]interface{}{
		"mfas": mfasValid,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "authenInfoRepository.UpdateOneByAttr(")
	}
	return &pb.MessageResponse{
		Ok:      true,
		Message: "Update mfa successful",
	}, nil

}
