package mongoImpl

import (
	"context"
	"fmt"
	"time"

	"github.com/thteam47/common-libs/errutil"
	"github.com/thteam47/common-libs/util"
	mongorepo "github.com/thteam47/go-survey-api/pkg/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryImpl struct {
	mongoDB *mongo.Collection
	timeOut time.Duration
}

func NewMongoRepo(mongoDB *mongo.Collection, timeOut time.Duration) mongorepo.MongoRepository {
	return &MongoRepositoryImpl{
		mongoDB: mongoDB,
		timeOut: timeOut,
	}
}

func (inst *MongoRepositoryImpl) InsertOne(data interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.timeOut)
	defer cancel()
	result, err := inst.mongoDB.InsertOne(ctx, data)
	if err != nil {
		return "", errutil.Wrap(err, "MongoDB.InsertOne")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
func (inst *MongoRepositoryImpl) GetOneByAttr(data map[string]string, result interface{}) error {
	if data == nil {
		return nil
	}
	findquery := []bson.M{}
	for _, key := range util.Keys[string, string](data) {
		value := ""
		if item, ok := data[key]; ok {
			value = item
			if key == "_id" {
				id, err := primitive.ObjectIDFromHex(value)
				if err != nil {
					return errutil.Wrap(err, "primitive.ObjectIDFromHex")
				}
				findquery = append(findquery, bson.M{
					key: id,
				})
				continue
			}
			findquery = append(findquery, bson.M{
				key: value,
			})
		}
	}

	filters := bson.M{
		"$or": findquery,
	}
	fmt.Println(data)

	err := inst.mongoDB.FindOne(context.Background(), filters).Decode(result)
	fmt.Println(result)
	if err != nil {
		return errutil.Wrap(err, "MongoDB.FindOne")
	}
	return nil
}
func (inst *MongoRepositoryImpl) UpdatebyId(data interface{}, id string) error {
	idPri, _ := primitive.ObjectIDFromHex(id)

	filterUser := bson.M{"_id": idPri}

	pByte, err := bson.Marshal(data)
	if err != nil {
		return errutil.Wrap(err, "bson.Marshal")
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return errutil.Wrap(err, "bson.Unmarshal")
	}
	_, err = inst.mongoDB.ReplaceOne(context.Background(), filterUser, update)
	if err != nil {
		return errutil.Wrap(err, "MongoDB.ReplaceOne")
	}
	return nil
}
func (inst *MongoRepositoryImpl) DeleteById(id string) error {
	idPri, _ := primitive.ObjectIDFromHex(id)
	_, err := inst.mongoDB.DeleteOne(context.Background(), bson.M{"_id": idPri})

	if err != nil {
		return errutil.Wrap(err, "MongoDB.DeleteOne")
	}
	return nil
}

func (inst *MongoRepositoryImpl) UpdateOneByAttr(userId string, data map[string]interface{}) error {
	filterUser := bson.M{"user_id": userId}
	dataUpdate := bson.M{}
	for _, key := range util.Keys[string, interface{}](data) {
		if value, ok := data[key]; ok {
			dataUpdate[key] = value
		}
	}
	opts := options.Update().SetUpsert(true)

	updateUser := bson.M{"$set": dataUpdate}
	_, err := inst.mongoDB.UpdateOne(context.Background(), filterUser, updateUser, opts)
	if err != nil {
		return errutil.Wrapf(err, "MongoDB.UpdateOne")
	}
	return nil
}

func (inst *MongoRepositoryImpl) Find(filter interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	cur, err := inst.mongoDB.Find(context.Background(), filter, options)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.Find")
	}
	return cur, nil
}
