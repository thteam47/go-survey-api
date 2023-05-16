package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/mongoutil"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/pkg/mongorepository"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type SurveyRepository struct {
	config         *SurveyRepositoryConfig
	baseRepository *mongorepository.BaseRepository
}

type SurveyRepositoryConfig struct {
	MongoClientWrapper *mongoutil.MongoClientWrapper `mapstructure:"mongo-client-wrapper"`
}

func NewSurveyRepositoryWithConfig(properties confg.Confg) (*SurveyRepository, error) {
	config := SurveyRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	mongoClientWrapper, err := mongoutil.NewBaseMongoClientWrapperWithConfig(properties.Sub("mongo-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewBaseMongoClientWrapperWithConfig")
	}
	return NewSurveyRepository(&SurveyRepositoryConfig{
		MongoClientWrapper: mongoClientWrapper,
	})
}

func NewSurveyRepository(config *SurveyRepositoryConfig) (*SurveyRepository, error) {
	inst := &SurveyRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = mongorepository.NewBaseRepository(&mongorepository.BaseRepositoryConfig{
		MongoClientWrapper: inst.config.MongoClientWrapper,
		Prototype:          models.Survey{},
		MongoIdField:       "Id",
		IdField:            "SurveyId",
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *SurveyRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Survey, error) {
	result := []models.Survey{}
	err := inst.baseRepository.FindAll(userContext, findRequest, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return result, nil
}

func (inst *SurveyRepository) Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error) {
	result, err := inst.baseRepository.Count(userContext, findRequest, &mongorepository.FindOptions{})
	if err != nil {
		return 0, errutil.Wrap(err, "baseRepository.Count")
	}
	return int32(result), nil
}

func (inst *SurveyRepository) FindById(userContext entity.UserContext, id string) (*models.Survey, error) {
	result := &models.Survey{}
	err := inst.baseRepository.FindOneByAttribute(userContext, "SurveyId", id, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *SurveyRepository) Create(userContext entity.UserContext, data *models.Survey) (*models.Survey, error) {
	err := inst.baseRepository.Create(userContext, data, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.Create")
	}
	return data, nil
}

func (inst *SurveyRepository) Update(userContext entity.UserContext, data *models.Survey, updateRequest *entity.UpdateRequest) (*models.Survey, error) {
	excludedProperties := []string{
		"CreatedTime", "Type",
	}

	err := inst.baseRepository.UpdateOneByAttribute(userContext, "SurveyId", data.SurveyId, data, updateRequest, &mongorepository.UpdateOptions{
		ExcludedProperties: excludedProperties,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.UpdateOneByAttribute")
	}
	return data, nil
}
func (inst *SurveyRepository) DeleteById(userContext entity.UserContext, id string) error {
	err := inst.baseRepository.DeleteOneByAttribute(userContext, "SurveyId", id)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteOneByAttribute")
	}
	return nil
}
