package servergrpc

import (
	"net"

	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-survey-api/pkg/grpcapp"
	"github.com/thteam47/go-survey-api/pkg/pb"
	repoimpl "github.com/thteam47/go-survey-api/pkg/repository/default"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(lis net.Listener, handler *handler.Handler) error {
	authRepository := grpcauth.NewAuthInterceptor(handler)
	// grpcConfig, err := grpcclient.NewGrpcClientConnWithConfig(handler.GrpcConnConfig)
	// if err != nil {
	// 	return errutil.Wrapf(err, "utilgrpc.NewGrpcClientConnWithConfig")
	// }
	// userService, err := repoimpl.NewUserRepo(grpcConfig)
	// if err != nil {
	// 	return errutil.Wrapf(err, "repoimpl.NewUserRepo")
	// }
	surveyRepository := repoimpl.NewSurveyRepo(handler.MongoRepository)

	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)
	pb.RegisterSurveyServiceServer(s, grpcapp.NewSurveyService(authRepository, surveyRepository))
	reflection.Register(s)
	return s.Serve(lis)
}
