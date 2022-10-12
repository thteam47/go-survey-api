package repoimpl

import (
	"context"
	"strings"

	"github.com/thteam47/go-identity-authen-api/errutil"
	v1 "github.com/thteam47/go-identity-authen-api/pkg/api-client/identity-api"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	"github.com/thteam47/go-identity-authen-api/pkg/repository"
	"github.com/thteam47/go-identity-authen-api/util"
	"google.golang.org/grpc"
)

type UserRepositoryImpl struct {
	config *util.GrpcClientConn
	client v1.IdentityServiceClient
}

func NewUserRepo(config *util.GrpcClientConn) (repository.UserRepository, error) {
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		return &UserRepositoryImpl{}, errutil.Wrapf(err, "grpc.Dial")
	}
	client := v1.NewIdentityServiceClient(conn)
	return &UserRepositoryImpl{
		config: config,
		client: client,
	}, nil
}

func (inst *UserRepositoryImpl) requestCtx() *v1.Context {
	return &v1.Context{
		AccessToken: inst.config.Config.AccessToken,
	}
}

func getUser(item *v1.User) (*models.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &models.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getUsers(items []*v1.User) ([]*models.User, error) {
	users := []*models.User{}
	for _, item := range items {
		user, err := getUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func makeUser(item *models.User) (*v1.User, error) {
	user := &v1.User{}
	err := util.ToMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return user, nil
}

func makeUsers(items []*models.User) ([]*v1.User, error) {
	users := []*v1.User{}
	for _, item := range items {
		user, err := makeUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func (inst *UserRepositoryImpl) FindByLoginName(userContext grpcauth.UserContext, name string) (*v1.User, error) {
	user := &v1.User{}
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Config.Timeout)
	defer cancel()
	if strings.Contains(name, "@") {
		user, err = inst.client.GetByEmail(ctx, &v1.StringRequest{
			Ctx:   inst.requestCtx(),
			Value: name,
		})
		if err != nil {
			return nil, errutil.Wrapf(err, "client.GetByEmail")
		}
	} else {
		user, err = inst.client.GetByLoginName(context.Background(), &v1.StringRequest{
			Ctx:   inst.requestCtx(),
			Value: name,
		})
		if err != nil {
			return nil, errutil.Wrapf(err, "client.GetByLoginName")
		}
	}
	return user, nil
}

func (inst *UserRepositoryImpl) FindById(userContext grpcauth.UserContext, id string) (*v1.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Config.Timeout)
	defer cancel()
	user, err := inst.client.GetById(ctx, &v1.StringRequest{
		Ctx:   inst.requestCtx(),
		Value: id,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetById")
	}

	return user, nil
}

func (inst *UserRepositoryImpl) Create(userContext grpcauth.UserContext, user *v1.User) (*v1.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Config.Timeout)
	defer cancel()
	result, err := inst.client.Create(ctx, &v1.UserRequest{
		Ctx:  inst.requestCtx(),
		Data: user,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetById")
	}

	return result, nil
}
func (inst *UserRepositoryImpl) GetAll(userContext grpcauth.UserContext, number int32, limit int32) ([]*models.User, error) {

	return nil, nil
}

func (inst *UserRepositoryImpl) Count(userContext grpcauth.UserContext) (int32, error) {

	return 0, nil
}

func (inst *UserRepositoryImpl) GetOneByAttr(userContext grpcauth.UserContext, data map[string]string) (*models.User, error) {

	return nil, nil
}

func (inst *UserRepositoryImpl) UpdatebyId(userContext grpcauth.UserContext, user *models.User, id string) (*models.User, error) {

	return nil, nil
}

func (inst *UserRepositoryImpl) DeleteById(userContext grpcauth.UserContext, id string) error {

	return nil
}

func (inst *UserRepositoryImpl) VerifyUser(userContext  grpcauth.UserContext, id string) error {
	userId := id
	if id == "" {
		userId = userContext.GetUserId()
	}
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Config.Timeout)
	defer cancel()
	_, err := inst.client.ApproveUser(ctx, &v1.ApproveUserRequest{
		Ctx:  inst.requestCtx(),
		UserId: userId,
		Status: "verified",
	})
	if err != nil {
		return errutil.Wrapf(err, "client.GetById")
	}

	return nil
}
