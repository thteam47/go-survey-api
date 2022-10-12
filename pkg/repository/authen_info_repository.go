package repository

import (
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
)

type AuthenInfoRepository interface {
	// GetAll(userContext grpcauth.UserContext, number int32, limit int32) ([]*models.User, error)
	// Count(userContext grpcauth.UserContext) (int32, error)
	GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.AuthenInfo, error)
	Create(userContext grpcauth.UserContext, item *models.AuthenInfo) (*models.AuthenInfo, error)
	UpdateOneByAttr(userContext grpcauth.UserContext, userId string, data map[string]interface{}) error
	DeleteOneByUserId(userContext grpcauth.UserContext, id string) error
	ForgotPassword(userContext grpcauth.UserContext, data string) (string, error)
	RegisterUser(userContext grpcauth.UserContext, username string, fullName string, email string) (string, error)
}
