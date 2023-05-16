package defaultcomponent

import (
	"context"
	"time"

	"github.com/thteam47/common-libs/confg"
	v1 "github.com/thteam47/common/api/recommend-api"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
	"github.com/thteam47/go-survey-api/util"
	"google.golang.org/grpc"
)

type RecommendService struct {
	config *RecommendServiceConfig
	client v1.RecommendServiceClient
}

type RecommendServiceConfig struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	AccessToken string        `mapstructure:"access_token"`
}

func NewRecommendServiceWithConfig(properties confg.Confg) (*RecommendService, error) {
	config := RecommendServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewRecommendService(&config)
}

func NewRecommendService(config *RecommendServiceConfig) (*RecommendService, error) {
	inst := &RecommendService{
		config: config,
	}
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		return nil, errutil.Wrapf(err, "grpc.Dial")
	}
	client := v1.NewRecommendServiceClient(conn)
	inst.client = client
	return inst, nil
}

func (inst *RecommendService) requestCtx(userContext entity.UserContext) *v1.Context {
	return &v1.Context{
		AccessToken: inst.config.AccessToken,
		DomainId:    userContext.DomainId(),
	}
}

func getCombinedData(item *v1.CombinedData) (*models.CombinedData, error) {
	if item == nil {
		return nil, nil
	}
	combinedData := &models.CombinedData{}
	err := util.FromMessage(item, combinedData)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return combinedData, nil
}

func getCombinedDatas(items []*v1.CombinedData) ([]models.CombinedData, error) {
	combinedDatas := []models.CombinedData{}
	for _, item := range items {
		combinedData, err := getCombinedData(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getCombinedData")
		}
		combinedDatas = append(combinedDatas, *combinedData)
	}
	return combinedDatas, nil
}

func (inst *RecommendService) GetCombinedDatasByDomain(userContext entity.UserContext, tenantId string) (*models.CombinedData, error) {
	result, err := inst.client.CombinedDataGetByTenantId(context.Background(), &v1.StringRequest{
		Ctx:   inst.requestCtx(userContext),
		Value: tenantId,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetById")
	}

	combinedDataItem, err := getCombinedData(result.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getCombinedDatas")
	}
	return combinedDataItem, nil
}
