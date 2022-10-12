package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang/glog"
	"github.com/soheilhy/cmux"
	"github.com/thteam47/common/handler"
	clienthttp "github.com/thteam47/go-survey-api/cmd/client"
	servergrpc "github.com/thteam47/go-survey-api/cmd/server"
	"github.com/thteam47/go-survey-api/pkg/configs"
)

func client(lis net.Listener, grpc_port string, http_port string) error {
	flag.Parse()
	defer glog.Flush()
	return clienthttp.Run(lis, grpc_port, http_port)
}

func serverGrpc(lis net.Listener, handler *handler.Handler) error {
	flag.Parse()
	defer glog.Flush()
	return servergrpc.Run(lis, handler)
}

func main() {
	cf, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	handler, err := handler.NewHandlerWithConfig(cf)
	if err != nil {
		log.Fatalln("NewHandlerWithConfig", err)
	}

	lis, err := net.Listen("tcp", cf.GrpcPort)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Server run on", cf.GrpcPort)

	m := cmux.New(lis)
	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())
	go serverGrpc(grpcL, handler)
	go client(httpL, cf.GrpcPort, cf.HttpPort)
	m.Serve()
}
