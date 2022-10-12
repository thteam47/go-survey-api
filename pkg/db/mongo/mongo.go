package mongorepo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	InsertOne(data interface{}) (string, error)
	GetOneByAttr(data map[string]string, result interface{}) error
	UpdatebyId(data interface{}, id string) error
	UpdateOneByAttr(userId string, data map[string]interface{}) error
	DeleteById(id string) error
	Find(filter interface{}, options *options.FindOptions) (cur *mongo.Cursor, err error)
}
