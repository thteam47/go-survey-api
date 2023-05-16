package grpcapp

import (
	"context"
	"fmt"

	"github.com/thteam47/common/entity"

	pb "github.com/thteam47/common/api/survey-api"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/pkg/adapter"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/component"
	"github.com/thteam47/go-survey-api/pkg/models"
	"github.com/thteam47/go-survey-api/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SurveyService struct {
	pb.SurveyServiceServer
	componentsContainer *component.ComponentsContainer
}

func NewSurveyService(componentsContainer *component.ComponentsContainer) *SurveyService {
	return &SurveyService{
		componentsContainer: componentsContainer,
	}
}

func getSurvey(item *pb.Survey) (*models.Survey, error) {
	if item == nil {
		return nil, nil
	}
	survey := &models.Survey{}
	err := util.FromMessage(item, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return survey, nil
}

func getSurveys(items []*pb.Survey) ([]*models.Survey, error) {
	surveys := []*models.Survey{}
	for _, item := range items {
		survey, err := getSurvey(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getSurveys")
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}

func makeSurvey(item *models.Survey) (*pb.Survey, error) {
	if item == nil {
		return nil, nil
	}
	survey := &pb.Survey{}
	err := util.ToMessage(item, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	if survey != nil {
		survey.UserIdJoin = []string{}
		survey.UserIdCreate = ""
		survey.UserIdVerify = ""
	}
	return survey, nil
}

func makeSurveys(items []models.Survey) ([]*pb.Survey, error) {
	surveys := []*pb.Survey{}
	for _, item := range items {
		survey, err := makeSurvey(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeSurvey")
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}

func (inst *SurveyService) CreateSurvey(ctx context.Context, req *pb.SurveyRequest) (*pb.SurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	survey, err := getSurvey(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getSurvey")
	}
	result, err := inst.componentsContainer.SurveyRepository().Create(userContext, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "UserRepository.Create")
	}
	item, err := makeSurvey(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurvey")
	}
	return &pb.SurveyResponse{
		Data: item,
	}, nil
}
func (inst *SurveyService) UpdateSurveyById(ctx context.Context, req *pb.UpdateSurveyRequest) (*pb.SurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	survey, err := getSurvey(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getSurvey")
	}
	result, err := inst.componentsContainer.SurveyRepository().Update(userContext, survey, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.UpdatebyId")
	}
	item, err := makeSurvey(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurvey")
	}
	return &pb.SurveyResponse{
		Data: item,
	}, nil
}
func (inst *SurveyService) GetAllSurvey(ctx context.Context, req *pb.ListRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Fobbiden!")
	// }
	findRequest, err := adapter.GetFindRequest(req, req.RequestPayload)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprint(err))
	}
	result, err := inst.componentsContainer.SurveyRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	item, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	count, err := inst.componentsContainer.SurveyRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.Count")
	}

	return &pb.ListSurveyResponse{
		Data:  item,
		Total: count,
	}, nil
}
func (inst *SurveyService) GetSurveyById(ctx context.Context, req *pb.StringRequest) (*pb.SurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.SurveyRepository().FindById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindById")
	}
	item, err := makeSurvey(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurvey")
	}
	return &pb.SurveyResponse{
		Data: item,
	}, nil
}
func (inst *SurveyService) GetSurveyByUserJoin(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Fobbiden!")
	// }
	findRequest := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "UserIdJoin",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    req.Value,
			},
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    req.Ctx.DomainId,
			},
			entity.FindRequestFilter{
				Key:      "Status",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    "approved",
			},
		},
	}
	result, err := inst.componentsContainer.SurveyRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	items, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	for _, item := range items {
		item.UserIdJoin = []string{req.Value}
	}
	count, err := inst.componentsContainer.SurveyRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.Count")
	}

	return &pb.ListSurveyResponse{
		Data:  items,
		Total: count,
	}, nil
}
func (inst *SurveyService) GetSurveyByUserCreate(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Fobbiden!")
	// }
	findRequest := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "UserId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    req.Value,
			},
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    req.Ctx.DomainId,
			},
		},
	}
	result, err := inst.componentsContainer.SurveyRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	item, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	count, err := inst.componentsContainer.SurveyRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.Count")
	}

	return &pb.ListSurveyResponse{
		Data:  item,
		Total: count,
	}, nil
}
func (inst *SurveyService) ApproveBySurveyId(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	survey, err := inst.componentsContainer.SurveyRepository().FindById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindById")
	}
	if survey == nil {
		return nil, status.Errorf(codes.NotFound, "Survey not found")
	}
	if survey.Status == "approved" || survey.Status == "denied" {
		return nil, status.Errorf(codes.FailedPrecondition, "Survey approved or denied")
	}
	survey.Status = "approved"
	_, err = inst.componentsContainer.SurveyRepository().Update(userContext, survey, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.UpdatebyId")
	}
	return &pb.StringResponse{}, nil
}
func (inst *SurveyService) DeleteSurveyById(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "delete", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.componentsContainer.SurveyRepository().DeleteById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}

func (inst *SurveyService) GetSurveyByTenant(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:survey", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	// if !entityutil.ServiceOrAdminRole(userContext) {
	// 	return nil, status.Errorf(codes.PermissionDenied, "Fobbiden!")
	// }
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedDatasByDomain(userContext.Clone().EscalatePrivilege(), req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedDatasByDomain")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "combinedData not found")
	}
	findRequest := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    combinedData.TenantId1,
			},
			entity.FindRequestFilter{
				Key:      "Status",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    "approved",
			},
		},
	}
	result, err := inst.componentsContainer.SurveyRepository().FindAll(entity.NewUserContext("default"), findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	items, err := makeSurveys(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	findRequest2 := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    combinedData.TenantId2,
			},
			entity.FindRequestFilter{
				Key:      "Status",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    "approved",
			},
		},
	}
	result2, err := inst.componentsContainer.SurveyRepository().FindAll(entity.NewUserContext("default"), findRequest2)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	items2, err := makeSurveys(result2)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}

	response := append(items, items2...)

	return &pb.ListSurveyResponse{
		Data:  response,
		Total: int32(len(response)),
	}, nil
}
