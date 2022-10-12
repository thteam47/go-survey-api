package repoimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/thteam47/go-identity-authen-api/errutil"
	v1 "github.com/thteam47/go-identity-authen-api/pkg/api-client/identity-api"
	"github.com/thteam47/go-identity-authen-api/pkg/db"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	pb "github.com/thteam47/go-identity-authen-api/pkg/pb"
	"github.com/thteam47/go-identity-authen-api/pkg/repository"
	"github.com/thteam47/go-identity-authen-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthenInfoRepositoryImpl struct {
	handler       *db.Handler
	userService   repository.UserRepository
	jwtRepository repository.JwtRepository
}

func NewAuthenInfoRepo(handler *db.Handler, userService repository.UserRepository, jwtRepository repository.JwtRepository) repository.AuthenInfoRepository {
	return &AuthenInfoRepositoryImpl{
		handler:       handler,
		userService:   userService,
		jwtRepository: jwtRepository,
	}
}

func getAuthenInfo(item *pb.AuthenInfo) (*models.AuthenInfo, error) {
	if item == nil {
		return nil, nil
	}
	authenInfo := &models.AuthenInfo{}
	err := util.FromMessage(item, authenInfo)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return authenInfo, nil
}

func makeAuthenInfo(item *models.AuthenInfo) (*pb.AuthenInfo, error) {
	authenInfo := &pb.AuthenInfo{}
	err := util.ToMessage(item, authenInfo)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return authenInfo, nil
}
func (inst *AuthenInfoRepositoryImpl) Create(userContext grpcauth.UserContext, item *models.AuthenInfo) (*models.AuthenInfo, error) {
	item.CreateTime = int32(time.Now().Unix())
	result, err := inst.handler.MongoDB.InsertOne(context.Background(), item)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.InsertOne")
	}
	uid := result.InsertedID.(primitive.ObjectID)
	item.ID = uid
	return item, nil
}

func (inst *AuthenInfoRepositoryImpl) GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.AuthenInfo, error) {
	if data == nil {
		return nil, nil
	}
	findquery := []bson.M{}
	for _, key := range util.Keys[string, string](data) {
		value := ""
		if item, ok := data[key]; ok {
			value = item
		}
		if key == "_id" {
			id, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				return nil, errutil.Wrap(err, "primitive.ObjectIDFromHex")
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

	filters := bson.M{
		"$or": findquery,
	}

	item := &models.AuthenInfo{}
	err := inst.handler.MongoDB.FindOne(context.Background(), filters).Decode(item)
	if err != nil {
		return nil, errutil.Wrap(err, "MongoDB.FindOne")
	}
	return item, nil
}

// func (inst *AuthenInfoRepositoryImpl) ChangeActionUser(idUser string, role string, a []string) error {
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
// 	_, err := inst.MongoDB.Collection(vi.GetString("collectionUser")).UpdateOne(ctx, filterUser, updateUser)
// 	if err != nil {
// 		return err
// 	}

//		return nil
//	}
func (inst *AuthenInfoRepositoryImpl) UpdateOneByAttr(userContext grpcauth.UserContext, userId string, data map[string]interface{}) error {
	filterUser := bson.M{"user_id": userId}
	dataUpdate := bson.M{}
	for _, key := range util.Keys[string, interface{}](data) {
		if value, ok := data[key]; ok {
			dataUpdate[key] = value
		}
	}
	opts := options.Update().SetUpsert(true)

	updateUser := bson.M{"$set": dataUpdate}
	_, err := inst.handler.MongoDB.UpdateOne(context.Background(), filterUser, updateUser, opts)
	if err != nil {
		return errutil.Wrapf(err, "MongoDB.UpdateOne")
	}
	return nil
}

// func (u *AuthenInfoRepositoryImpl) ChangePassUser(idUser string, pass string) error {
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
func (inst *AuthenInfoRepositoryImpl) DeleteOneByUserId(userContext grpcauth.UserContext, id string) error {
	_, err := inst.handler.MongoDB.DeleteOne(context.Background(), bson.M{"user_id": id})

	if err != nil {
		return errutil.Wrap(err, "MongoDB.DeleteOne")
	}
	return nil
}

func (inst *AuthenInfoRepositoryImpl) ForgotPassword(userContext grpcauth.UserContext, data string) (string, error) {
	userItem, err := inst.userService.FindByLoginName(userContext, data)
	if err != nil {
		return "", errutil.Wrapf(err, "userService.FindByLoginName")
	}
	if userItem == nil {
		return "", errutil.NewWithMessage("Username or password incorrect")
	}
	tokenInfo := &models.TokenInfo{
		AuthenticationDone: true,
		UserId:             userItem.UserId,
		Exp:                int32(time.Now().Add(5 * time.Minute).Unix()),
	}
	token, err := inst.jwtRepository.Generate(tokenInfo)
	if err != nil {
		return "", errutil.Wrapf(err, "jwtRepository.Generate")
	}
	dataMail, err := util.ParseTemplate("../util/template.html", map[string]string{
		"message":    "Click the link to change password.",
		"username":   userItem.FullName,
		"title":      "Forgot Password",
		"buttonText": "Change Passowrd Now",
		"link":       fmt.Sprintf("http://localhost:4200/update-password/%s", token),
	})
	if err != nil {
		return "", errutil.Wrapf(err, "util.ParseTemplate")
	}
	err = util.SendMail([]string{userItem.Email}, dataMail)
	if err != nil {
		return "", errutil.Wrapf(err, "util.SendMail")
	}
    return fmt.Sprintf("Click the link in your email %s to change your password", userItem.Email), nil
}

func (inst *AuthenInfoRepositoryImpl) RegisterUser(userContext grpcauth.UserContext, username string, fullName string, email string) (string, error) {
	userData := &v1.User{
		FullName:   fullName,
		Email:      email,
		Username:   username,
		Role:       "member",
		CreateTime: int32(time.Now().Unix()),
		Status:     "pending",
	}
	userItem, err := inst.userService.Create(userContext, userData)
	if err != nil {
		return "", errutil.Wrapf(err, "userService.FindByLoginName")
	}
	if userItem == nil {
		return "", errutil.NewWithMessage("Username or password incorrect")
	}
	tokenInfo := &models.TokenInfo{
		AuthenticationDone: true,
		UserId:             userItem.UserId,
		Exp:                int32(time.Now().Add(5 * time.Minute).Unix()),
	}
	token, err := inst.jwtRepository.Generate(tokenInfo)
	if err != nil {
		return "", errutil.Wrapf(err, "jwtRepository.Generate")
	}
	dataMail, err := util.ParseTemplate("../util/template.html", map[string]string{
		"message":    "Click the link to verify account.",
		"username":   fullName,
		"title":      "Verify Account",
		"buttonText": "Verify Now",
		"link":       fmt.Sprintf("http://localhost:4200/verify-account/%s", token),
	})
	if err != nil {
		return "", errutil.Wrapf(err, "util.ParseTemplate")
	}
	err = util.SendMail([]string{userItem.Email}, dataMail)
	if err != nil {
		return "", errutil.Wrapf(err, "util.SendMail")
	}
	return fmt.Sprintf("Click the link in your email %s to verify your account", userData.Email), nil
}
