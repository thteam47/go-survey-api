package servergrpc

import (
	"net"

	"github.com/thteam47/go-identity-authen-api/errutil"
	"github.com/thteam47/go-identity-authen-api/pkg/db"
	"github.com/thteam47/go-identity-authen-api/pkg/grpcapp"
	grpcauth "github.com/thteam47/go-identity-authen-api/pkg/grpcutil"
	"github.com/thteam47/go-identity-authen-api/pkg/pb"
	"github.com/thteam47/go-identity-authen-api/pkg/repository"
	repoimpl "github.com/thteam47/go-identity-authen-api/pkg/repository/default"
	"github.com/thteam47/go-identity-authen-api/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(lis net.Listener, handler *db.Handler) error {
	authRepository := grpcauth.NewAuthInterceptor(handler)
	grpcConfig, err := util.NewGrpcClientConnWithConfig(handler.GrpcConnConfig)
	if err != nil {
		return errutil.Wrapf(err, "util.NewGrpcClientConnWithConfig")
	}
	userService, err := repoimpl.NewUserRepo(grpcConfig)
	if err != nil {
		return errutil.Wrapf(err, "repoimpl.NewUserRepo")
	}
	jwtRepository := repoimpl.NewJwtRepo(handler)
	if err != nil {
		return errutil.Wrapf(err, "repoimpl.NewUserRepo")
	}
	authenInfoRepository := repoimpl.NewAuthenInfoRepo(handler, userService, jwtRepository)
	if err != nil {
		return errutil.Wrapf(err, "repoimpl.NewUserRepo")
	}
	componentContanier := repository.NewComponentContanier(userService, authenInfoRepository, jwtRepository, authRepository, handler)
	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)
	pb.RegisterIdentityAuthenServiceServer(s, grpcapp.NewIdentityAuthenService(handler, componentContanier, userService, authRepository, authenInfoRepository))
	reflection.Register(s)
	return s.Serve(lis)
}
