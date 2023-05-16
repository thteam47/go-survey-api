package component

import grpcauth "github.com/thteam47/common/grpcutil"

type ComponentFactory interface {
	CreateAuthService() *grpcauth.AuthInterceptor
	CreateCategoryRepository() (CategoryRepository, error)
	CreateSurveyRepository() (SurveyRepository, error)
	CreateIdentityService() (IdentityService, error)
	CreateRecommendService() (RecommendService, error)
}
