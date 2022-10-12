package repository

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type SurveyRepository interface {
	GetAll(userContext grpcauth.UserContext, number int32, limit int32, filter interface{}) ([]*models.Survey, error)
	Count(userContext grpcauth.UserContext) (int32, error)
	GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.Survey, error)
	Create(userContext grpcauth.UserContext, Survey *models.Survey) (*models.Survey, error)
	UpdatebyId(userContext grpcauth.UserContext, Survey *models.Survey, id string) (*models.Survey, error)
	DeleteById(userContext grpcauth.UserContext, id string) error
	ApproveBySurveyId(userContext grpcauth.UserContext, id string) error
}
