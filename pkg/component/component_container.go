package component

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/db"
)

type ComponentsContainer struct {
	categoryRepository CategoryRepository
	authService        *grpcauth.AuthInterceptor
	handler            *db.Handler
	surveyRepository   SurveyRepository
	identityService    IdentityService
	recommendService   RecommendService
}

func NewComponentsContainer(componentFactory ComponentFactory) (*ComponentsContainer, error) {
	inst := &ComponentsContainer{}

	var err error
	inst.authService = componentFactory.CreateAuthService()
	inst.categoryRepository, err = componentFactory.CreateCategoryRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateAuthenInfoRepository")
	}
	inst.surveyRepository, err = componentFactory.CreateSurveyRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateSurveyRepository")
	}
	inst.identityService, err = componentFactory.CreateIdentityService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateIdentityService")
	}
	inst.recommendService, err = componentFactory.CreateRecommendService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreatereRommendService")
	}
	return inst, nil
}

func (inst *ComponentsContainer) AuthService() *grpcauth.AuthInterceptor {
	return inst.authService
}

func (inst *ComponentsContainer) CategoryRepository() CategoryRepository {
	return inst.categoryRepository
}

func (inst *ComponentsContainer) SurveyRepository() SurveyRepository {
	return inst.surveyRepository
}

func (inst *ComponentsContainer) IdentityService() IdentityService {
	return inst.identityService
}

func (inst *ComponentsContainer) RecommendService() RecommendService {
	return inst.recommendService
}

var errorCodeBadRequest = 400
