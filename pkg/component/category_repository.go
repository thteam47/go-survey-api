package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type CategoryRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Category, error)
	Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error)
	FindById(userContext entity.UserContext, id string) (*models.Category, error)
	Create(userContext entity.UserContext, user *models.Category) (*models.Category, error)
	Update(userContext entity.UserContext, data *models.Category, updateRequest *entity.UpdateRequest) (*models.Category, error)
	DeleteById(userContext entity.UserContext, id string) error
}
