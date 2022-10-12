package grpcapp

import (
	"context"
	"fmt"

	"github.com/thteam47/common-libs/errutil"
	"github.com/thteam47/common-libs/util"
	grpcauth "github.com/thteam47/common/grpcutil"
	pbCommon "github.com/thteam47/common/pb"
	"github.com/thteam47/go-survey-api/pkg/models"
	"github.com/thteam47/go-survey-api/pkg/pb"
	"github.com/thteam47/go-survey-api/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SurveyService struct {
	pb.SurveyServiceServer
	authRepository   *grpcauth.AuthInterceptor
	surveyRepository repository.SurveyRepository
}

func NewSurveyService(authRepository *grpcauth.AuthInterceptor, surveyRepository repository.SurveyRepository) *SurveyService {
	return &SurveyService{
		authRepository:   authRepository,
		surveyRepository: surveyRepository,
	}
}

func getSurvey(item *pb.Survey) (*models.Survey, error) {
	if item == nil {
		return nil, nil
	}
	Survey := &models.Survey{}
	err := util.FromMessage(item, Survey)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return Survey, nil
}

func getSurveys(items []*pb.Survey) ([]*models.Survey, error) {
	Surveys := []*models.Survey{}
	for _, item := range items {
		Survey, err := getSurvey(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getSurvey")
		}
		Surveys = append(Surveys, Survey)
	}
	return Surveys, nil
}

func makeSurvey(item *models.Survey) (*pb.Survey, error) {
	Survey := &pb.Survey{}
	err := util.ToMessage(item, Survey)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return Survey, nil
}

func makeSurveys(items []*models.Survey) ([]*pb.Survey, error) {
	Surveys := []*pb.Survey{}
	for _, item := range items {
		Survey, err := makeSurvey(item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeSurvey")
		}
		Surveys = append(Surveys, Survey)
	}
	return Surveys, nil
}

func (inst *SurveyService) CreateSurvey(ctx context.Context, req *pb.SurveyRequest) (*pb.Survey, error) {
	userContext, err := inst.authRepository.Authentication(ctx, &pbCommon.Context{
		AccessToken: req.Ctx.AccessToken,
	}, "@any", "@any")
	fmt.Println("fdsghsdh")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	survey, err := getSurvey(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getSurvey")
	}
	survey.UserIdCreate = userContext.GetUserId()

	result, err := inst.surveyRepository.Create(nil, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.Create")
	}
	item, err := makeSurvey(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurvey")
	}
	return item, nil
}

func (inst *SurveyService) GetSurveyById(ctx context.Context, req *pb.StringRequest) (*pb.Survey, error) {
	// SurveyContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "Survey-api:Survey", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	fmt.Println(req)
	result, err := inst.surveyRepository.GetOneByAttr(nil, map[string]string{
		"_id": req.Value,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.GetById")
	}
	item, err := makeSurvey(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurvey")
	}
	return item, nil
}

// func (inst *SurveyService) GetByEmail(ctx context.Context, req *pb.StringRequest) (*pb.Survey, error) {
// 	SurveyContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "Survey-api:Survey", "get")
// 	if err != nil {
// 		return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
// 	}
// 	result, err := inst.SurveyRepository.GetOneByAttr(nil, map[string]string{
// 		"email": req.Value,
// 	})
// 	if err != nil {
// 		return nil, errutil.Wrap(err, "SurveyRepository.GetById")
// 	}
// 	item, err := makeSurvey(result)
// 	if err != nil {
// 		return nil, errutil.Wrap(err, "makeSurvey")
// 	}
// 	return item, nil
// }

func (inst *SurveyService) GetAll(ctx context.Context, req *pb.ListRequest) (*pb.ListSurveyResponse, error) {
	// SurveyContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "Survey-api:Survey", "get")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	number := 1
	limit := 10
	if req != nil && req.Data != nil {
		if req.Data.Limit > 0 {
			limit = int(req.Data.Limit)
		}
		if req.Data.Number >= 1 {
			number = int(req.Data.Number)
		}
	}
	result, err := inst.surveyRepository.GetAll(nil, int32(number), int32(limit), bson.M{})
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.GetAll")
	}
	item, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	count, err := inst.surveyRepository.Count(nil)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.Count")
	}

	return &pb.ListSurveyResponse{
		Data:  item,
		Total: count,
	}, nil
}

func (inst *SurveyService) GetSurveyByUserCreate(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, &pbCommon.Context{
		AccessToken: req.Ctx.AccessToken,
	}, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	userIdCreate := req.Value
	if !userContext.ServiceOrHasRoleAdmin() {
		userIdCreate = userContext.GetUserId()
	}
	result, err := inst.surveyRepository.GetAll(nil, -1, -1, bson.M{
		"user_id_create": userIdCreate,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.GetAll")
	}
	item, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListSurveyResponse{
		Data:  item,
		Total: int32(len(item)),
	}, nil
}

func (inst *SurveyService) GetSurveyByUserJoin(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, &pbCommon.Context{
		AccessToken: req.Ctx.AccessToken,
	}, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	userIdJoin := req.Value
	if !userContext.ServiceOrHasRoleAdmin() {
		userIdJoin = userContext.GetUserId()
	}
	result, err := inst.surveyRepository.GetAll(nil, -1, -1, bson.M{
		"user_id_join": userIdJoin,
		"status":       "approve",
	})
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.GetAll")
	}
	item, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListSurveyResponse{
		Data:  item,
		Total: int32(len(item)),
	}, nil
}

func (inst *SurveyService) ApproveBySurveyId(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.authRepository.Authentication(ctx, &pbCommon.Context{
		AccessToken: req.Ctx.AccessToken,
	}, "@any", "@any")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.surveyRepository.ApproveBySurveyId(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.GetAll")
	}
	return &pb.StringResponse{}, nil
}

func (inst *SurveyService) UpdateSurveyById(ctx context.Context, req *pb.UpdateSurveyRequest) (*pb.StringResponse, error) {
	// SurveyContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "Survey-api:Survey", "update")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	survey, err := getSurvey(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getSurvey")
	}
	_, err = inst.surveyRepository.UpdatebyId(nil, survey, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.UpdatebyId")
	}
	return &pb.StringResponse{}, nil
}

func (inst *SurveyService) DeleteSurveyById(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	// SurveyContext, err := inst.authRepository.Authentication(ctx, req.Ctx, "Survey-api:Survey", "delete")
	// if err != nil {
	// 	return nil, status.Errorf(codes.PermissionDenied, "authRepository.Authentication")
	// }
	err := inst.surveyRepository.DeleteById(nil, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}
