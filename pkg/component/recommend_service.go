package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type RecommendService interface {
	GetCombinedDatasByDomain(userContext entity.UserContext, tenantId string) (*models.CombinedData, error)
}
