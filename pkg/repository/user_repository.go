package repository

import (
	v1 "github.com/thteam47/go-identity-authen-api/pkg/api-client/identity-api"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
)

type UserRepository interface {
	FindByLoginName(userContext grpcauth.UserContext, name string) (*v1.User, error)
	FindById(userContext grpcauth.UserContext, id string) (*v1.User, error)
	Create(userContext grpcauth.UserContext, user *v1.User) (*v1.User, error)
	VerifyUser(userContext grpcauth.UserContext, id string) (error)
	// GetAll(userContext grpcauth.UserContext, number int32, limit int32) ([]*models.User, error)
	// Count(userContext grpcauth.UserContext) (int32, error)
	// GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.User, error)
	// Create(userContext grpcauth.UserContext, user *models.User) (*models.User, error)
	// UpdatebyId(userContext grpcauth.UserContext, user *models.User, id string) (*models.User, error)
	// DeleteById(userContext grpcauth.UserContext, id string) error
}
