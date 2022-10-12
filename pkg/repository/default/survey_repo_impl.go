package repoimpl

import (
	"context"
	"log"
	"time"

	"github.com/thteam47/common-libs/errutil"
	grpcauth "github.com/thteam47/common/grpcutil"
	mongorepo "github.com/thteam47/common/handler/mongo"
	"github.com/thteam47/go-survey-api/pkg/models"
	"github.com/thteam47/go-survey-api/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SurveyRepositoryImpl struct {
	mongoRepository mongorepo.MongoRepository
}

func NewSurveyRepo(mongoRepository mongorepo.MongoRepository) repository.SurveyRepository {
	return &SurveyRepositoryImpl{
		mongoRepository: mongoRepository,
	}
}

func (inst *SurveyRepositoryImpl) Create(userContext grpcauth.UserContext, survey *models.Survey) (*models.Survey, error) {
	survey.CreateTime = int32(time.Now().Unix())
	result, err := inst.mongoRepository.InsertOne(survey)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.InsertOne")
	}
	survey.SurveyId = result
	return survey, nil
}

func (inst *SurveyRepositoryImpl) GetAll(userContext grpcauth.UserContext, number int32, limit int32, filter interface{}) ([]*models.Survey, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	if number != -1 && limit != -1 {
		if number == 1 {
			findOptions.SetSkip(0)
			findOptions.SetLimit(int64(limit))
		} else {
			findOptions.SetSkip(int64((number - 1) * limit))
			findOptions.SetLimit(int64(limit))
		}
	}
	cur, err := inst.mongoRepository.Find(filter, findOptions)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.Find")
	}
	var surveys []*models.Survey
	for cur.Next(context.TODO()) {
		var elem models.Survey
		err = cur.Decode(&elem)
		elem.SurveyId = elem.ID.Hex()
		if err != nil {
			return nil, errutil.Wrap(err, "Decode")
		}
		surveys = append(surveys, &elem)
	}
	return surveys, nil
}

func (inst *SurveyRepositoryImpl) Count(userContext grpcauth.UserContext) (int32, error) {
	findOptions := options.Find()

	cur, err := inst.mongoRepository.Find(bson.M{}, findOptions)
	if err != nil {
		return 0, err
	}
	var surveys []*models.Survey
	for cur.Next(context.TODO()) {
		var elem *models.Survey
		er := cur.Decode(&elem)
		if er != nil {
			log.Fatal(err)
		}
		surveys = append(surveys, elem)
	}
	return int32(len(surveys)), nil
}

func (inst *SurveyRepositoryImpl) GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.Survey, error) {
	survey := &models.Survey{}
	err := inst.mongoRepository.GetOneByAttr(data, &survey)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.FindOne")
	}
	survey.SurveyId = survey.ID.Hex()
	return survey, nil
}

// func (inst *SurveyRepositoryImpl) ChangeActionUser(idUser string, role string, a []string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	var roleUser string
// 	if role == "" {
// 		roleUser = "staff"
// 	} else {
// 		roleUser = role
// 	}
// 	var actionList []string
// 	if roleUser == "admin" {
// 		actionList = append(actionList, "All Rights")
// 	} else if roleUser == "assistant" {
// 		actionList = []string{"Add Server", "Update Server", "Detail Status", "Export", "Connect", "Disconnect", "Delete Server", "Change Password"}
// 	} else {
// 		actionList = a
// 	}
// 	id, _ := primitive.ObjectIDFromHex(idUser)
// 	filterUser := bson.M{"_id": id}
// 	updateUser := bson.M{"$set": bson.M{
// 		"role":   roleUser,
// 		"action": actionList,
// 	}}
// 	_, err := inst.mongoRepository.Collection(vi.GetString("collectionUser")).UpdateOne(ctx, filterUser, updateUser)
// 	if err != nil {
// 		return err
// 	}

//		return nil
//	}
func (inst *SurveyRepositoryImpl) UpdatebyId(userContext grpcauth.UserContext, survey *models.Survey, id string) (*models.Survey, error) {
	survey.UpdateTime = int32(time.Now().Unix())

	err := inst.mongoRepository.UpdatebyId(survey, id)
	if err != nil {
		return nil, errutil.Wrap(err, "mongoRepository.UpdatebyId")
	}
	survey.SurveyId = id
	return survey, nil
}

func (inst *SurveyRepositoryImpl) ApproveBySurveyId(userContext grpcauth.UserContext, id string) error {
	updateTime := int32(time.Now().Unix())
	err := inst.mongoRepository.UpdateOneByAttr(id, map[string]interface{}{
		"update_time": updateTime,
		"status":      "approve",
	})
	if err != nil {
		return errutil.Wrap(err, "mongoRepository.UpdateOneByAttr")
	}
	return nil
}

// func (u *SurveyRepositoryImpl) ChangePassUser(idUser string, pass string) error {
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// var id primitive.ObjectID
// id, _ = primitive.ObjectIDFromHex(idUser)

// passHash, _ := drive.HashPassword(pass)
// filterUser := bson.M{"_id": id}
//
//	updateUser := bson.M{"$set": bson.M{
//		"password": passHash,
//	}}
//
// _, err := u.MongoDB.UpdateOne(ctx, filterUser, updateUser)
//
//	if err != nil {
//		return err
//	}
//
//		return nil
//	}
func (inst *SurveyRepositoryImpl) DeleteById(userContext grpcauth.UserContext, id string) error {
	err := inst.mongoRepository.DeleteById(id)

	if err != nil {
		return errutil.Wrap(err, "MongoDB.DeleteById")
	}
	return nil
}
