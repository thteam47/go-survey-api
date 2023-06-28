package grpcapp

import (
	"context"

	"github.com/thteam47/common/entity"

	"github.com/brianvoe/gofakeit/v6"
	pb "github.com/thteam47/common/api/survey-api"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *SurveyService) FakeSurveys(ctx context.Context, req *pb.FakeSurveyRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, http.StatusText(http.StatusForbidden))
	// }

	userIdJoin := []string{}
	users, _, err := inst.componentsContainer.IdentityService().GetUsersByDomain(userContext, req.Ctx.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "IdentityService.GetUsersByDomain")
	}
	for _, user := range users {
		userIdJoin = append(userIdJoin, user.UserId)
	}
	for i := 1; i <= int(req.NumberSurvey); i++ {
		categories, err := inst.componentsContainer.CategoryRepository().FindAll(userContext.Clone().EscalatePrivilege(), &entity.FindRequest{
			Limit: -1,
			Filters: []entity.FindRequestFilter{
				entity.FindRequestFilter{
					Key:      "DomainId",
					Operator: entity.FindRequestFilterOperatorEqualTo,
					Value:    req.Ctx.DomainId,
				},
			},
		})
		if err != nil {
			return nil, errutil.Wrap(err, "CategoryRepository.FindAll")
		}
		questions := []models.Question{}
		for _, category := range categories {
			answers := []string{}
			for j := 0; j < 5; j++ {
				answers = append(answers, gofakeit.Quote())
			}
			questions = append(questions, models.Question{
				CategoryId: category.CategoryId,
				Position:   category.Position,
				Type:       "radio",
				Message:    gofakeit.Question(),
				Answers:    answers,
			})
		}
		survey := &models.Survey{
			DomainId:   req.Ctx.DomainId,
			Name:       gofakeit.Sentence(20),
			Status:     "approved",
			UserIdJoin: userIdJoin,
			Questions:  questions,
		}

		_, err = inst.componentsContainer.SurveyRepository().Create(userContext, survey)
		if err != nil {
			return nil, errutil.Wrap(err, "SurveyRepository.Create")
		}
	}

	return &pb.StringResponse{}, nil
}
