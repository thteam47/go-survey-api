package grpcapp

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	pb "github.com/thteam47/common/api/survey-api"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *SurveyService) FakeCategories(ctx context.Context, req *pb.FakeCategoryRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, http.StatusText(http.StatusForbidden))
	// }
	for i := 1; i <= int(req.NumberCategory); i++ {
		category := &models.Category{
			DomainId: req.Ctx.DomainId,
			Name:     gofakeit.Fruit(),
			Position: int32(i),
		}
		_, err := inst.componentsContainer.CategoryRepository().Create(userContext, category)
		if err != nil {
			return nil, errutil.Wrap(err, "CategoryRepository.Create")
		}
	}

	return &pb.StringResponse{}, nil
}
