package component

import (
	"github.com/thteam47/common/entity"
)

type IdentityService interface {
	GetUsersByDomain(userContext entity.UserContext, tenantId string) ([]entity.User, int32, error)
}
