package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type SurveyRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Survey, error)
	Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error)
	FindById(userContext entity.UserContext, id string) (*models.Survey, error)
	Create(userContext entity.UserContext, user *models.Survey) (*models.Survey, error)
	Update(userContext entity.UserContext, data *models.Survey, updateRequest *entity.UpdateRequest) (*models.Survey, error)
	DeleteById(userContext entity.UserContext, id string) error
}
