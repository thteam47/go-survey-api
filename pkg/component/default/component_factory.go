package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/component"
)

type ComponentFactory struct {
	properties confg.Confg
	handle     *handler.Handler
}

func NewComponentFactory(properties confg.Confg, handle *handler.Handler) (*ComponentFactory, error) {
	inst := &ComponentFactory{
		properties: properties,
		handle:     handle,
	}

	return inst, nil
}

func (inst *ComponentFactory) CreateAuthService() *grpcauth.AuthInterceptor {
	authService := grpcauth.NewAuthInterceptor(inst.handle)
	return authService
}

func (inst *ComponentFactory) CreateCategoryRepository() (component.CategoryRepository, error) {
	categoryRepository, err := NewCategoryRepositoryWithConfig(inst.properties.Sub("category-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewCategoryRepositoryWithConfig")
	}
	return categoryRepository, nil
}

func (inst *ComponentFactory) CreateSurveyRepository() (component.SurveyRepository, error) {
	surveyRepository, err := NewSurveyRepositoryWithConfig(inst.properties.Sub("survey-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewSurveyRepositoryWithConfig")
	}
	return surveyRepository, nil
}

func (inst *ComponentFactory) CreateIdentityService() (component.IdentityService, error) {
	identityService, err := NewIdentityServiceWithConfig(inst.properties.Sub("identity-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewIdentityServiceWithConfig")
	}
	return identityService, nil
}

func (inst *ComponentFactory) CreateRecommendService() (component.RecommendService, error) {
	recommendService, err := NewRecommendServiceWithConfig(inst.properties.Sub("recommend-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewRecommendServiceWithConfig")
	}
	return recommendService, nil
}
