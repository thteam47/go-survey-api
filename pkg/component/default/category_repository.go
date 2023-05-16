package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/mongoutil"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/pkg/mongorepository"
	"github.com/thteam47/go-survey-api/errutil"
	"github.com/thteam47/go-survey-api/pkg/models"
)

type CategoryRepository struct {
	config         *CategoryRepositoryConfig
	baseRepository *mongorepository.BaseRepository
}

type CategoryRepositoryConfig struct {
	MongoClientWrapper *mongoutil.MongoClientWrapper `mapstructure:"mongo-client-wrapper"`
}

func NewCategoryRepositoryWithConfig(properties confg.Confg) (*CategoryRepository, error) {
	config := CategoryRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	mongoClientWrapper, err := mongoutil.NewBaseMongoClientWrapperWithConfig(properties.Sub("mongo-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewBaseMongoClientWrapperWithConfig")
	}
	return NewCategoryRepository(&CategoryRepositoryConfig{
		MongoClientWrapper: mongoClientWrapper,
	})
}

func NewCategoryRepository(config *CategoryRepositoryConfig) (*CategoryRepository, error) {
	inst := &CategoryRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = mongorepository.NewBaseRepository(&mongorepository.BaseRepositoryConfig{
		MongoClientWrapper: inst.config.MongoClientWrapper,
		Prototype:          models.Category{},
		MongoIdField:       "Id",
		IdField:            "CategoryId",
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *CategoryRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Category, error) {
	result := []models.Category{}
	err := inst.baseRepository.FindAll(userContext, findRequest, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return result, nil
}

func (inst *CategoryRepository) Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error) {
	result, err := inst.baseRepository.Count(userContext, findRequest, &mongorepository.FindOptions{})
	if err != nil {
		return 0, errutil.Wrap(err, "baseRepository.Count")
	}
	return int32(result), nil
}

func (inst *CategoryRepository) FindById(userContext entity.UserContext, id string) (*models.Category, error) {
	result := &models.Category{}
	err := inst.baseRepository.FindOneByAttribute(userContext, "CategoryId", id, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *CategoryRepository) Create(userContext entity.UserContext, data *models.Category) (*models.Category, error) {
	err := inst.baseRepository.Create(userContext, data, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.Create")
	}
	return data, nil
}

func (inst *CategoryRepository) Update(userContext entity.UserContext, data *models.Category, updateRequest *entity.UpdateRequest) (*models.Category, error) {
	excludedProperties := []string{
		"CreatedTime", "Type",
	}

	err := inst.baseRepository.UpdateOneByAttribute(userContext, "CategoryId", data.CategoryId, data, updateRequest, &mongorepository.UpdateOptions{
		ExcludedProperties: excludedProperties,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.UpdateOneByAttribute")
	}
	return data, nil
}
func (inst *CategoryRepository) DeleteById(userContext entity.UserContext, id string) error {
	err := inst.baseRepository.DeleteOneByAttribute(userContext, "CategoryId", id)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteOneByAttribute")
	}
	return nil
}
