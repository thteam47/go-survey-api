package grpcapp

import (
	"context"
	"fmt"

	pb "github.com/thteam47/common/api/survey-api"
	"github.com/thteam47/common/entity"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/pkg/adapter"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
	"github.com/thteam47/go-survey-api/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getCategory(item *pb.Category) (*models.Category, error) {
	if item == nil {
		return nil, nil
	}
	category := &models.Category{}
	err := util.FromMessage(item, category)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return category, nil
}

func getCategories(items []*pb.Category) ([]*models.Category, error) {
	categories := []*models.Category{}
	for _, item := range items {
		category, err := getCategory(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getCategories")
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func makeCategory(item *models.Category) (*pb.Category, error) {
	if item == nil {
		return nil, nil
	}
	category := &pb.Category{}
	err := util.ToMessage(item, category)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return category, nil
}

func makeCategories(items []models.Category) ([]*pb.Category, error) {
	categories := []*pb.Category{}
	for _, item := range items {
		category, err := makeCategory(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeCategory")
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (inst *SurveyService) CreateCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:category", "create", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	category, err := getCategory(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getCategory")
	}
	result, err := inst.componentsContainer.CategoryRepository().Create(userContext, category)
	if err != nil {
		return nil, errutil.Wrap(err, "UserRepository.Create")
	}
	item, err := makeCategory(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategory")
	}
	return &pb.CategoryResponse{
		Data: item,
	}, nil
}

func (inst *SurveyService) GetCategory(ctx context.Context, req *pb.StringRequest) (*pb.CategoryResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:category", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.CategoryRepository().FindById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindById")
	}
	item, err := makeCategory(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategory")
	}
	return &pb.CategoryResponse{
		Data: item,
	}, nil
}

func (inst *SurveyService) GetAllCategory(ctx context.Context, req *pb.ListRequest) (*pb.ListCategoryResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:category", "get", &grpcauth.AuthenOption{})
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
	result, err := inst.componentsContainer.CategoryRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.FindAll")
	}
	item, err := makeCategories(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategories")
	}
	count, err := inst.componentsContainer.CategoryRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.Count")
	}

	return &pb.ListCategoryResponse{
		Data:  item,
		Total: count,
	}, nil
}

func (inst *SurveyService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:category", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	category, err := getCategory(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getCategory")
	}
	result, err := inst.componentsContainer.CategoryRepository().Update(userContext, category, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.UpdatebyId")
	}
	item, err := makeCategory(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategory")
	}
	return &pb.CategoryResponse{
		Data: item,
	}, nil
}

func (inst *SurveyService) DeleteCategory(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "survey-api:category", "delete", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.componentsContainer.CategoryRepository().DeleteById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "userRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}

func (inst *SurveyService) GetCategoriesByRecommendTenant(ctx context.Context, req *pb.StringRequest) (*pb.ListCategoryResponse, error) {
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
		},
	}
	result, err := inst.componentsContainer.CategoryRepository().FindAll(entity.NewUserContext("default"), findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	items, err := makeCategories(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategories")
	}
	findRequest2 := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    combinedData.TenantId2,
			},
		},
	}
	result2, err := inst.componentsContainer.CategoryRepository().FindAll(entity.NewUserContext("default"), findRequest2)
	if err != nil {
		return nil, errutil.Wrap(err, "surveyRepository.FindAll")
	}
	items2, err := makeCategories(result2)
	if err != nil {
		return nil, errutil.Wrap(err, "makeCategories")
	}

	lenItems := len(items)
	for key, value := range items2 {
		value.Position = int32(lenItems + key + 1)
		items = append(items, value)
	}

	return &pb.ListCategoryResponse{
		Data:  items,
		Total: int32(len(items)),
	}, nil
}
