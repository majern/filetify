package server

import (
	"fmt"
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func StartServer() {
	shared.InitLogger(&GetConfiguration().LogConfig)
	logrus.Info("Starting server")

	startGrpcServer()

	select {}
}

func startGrpcServer() {
	endpoint := fmt.Sprintf("localhost:%d", GetConfiguration().Rcp.Port)
	lis, err := net.Listen("tcp", endpoint)
	shared.HandleErrorWithMsg(err, true, "Failed to start gRPC server")
	grpcServer := grpc.NewServer()
	proto_v1.RegisterSynchronizationServiceServer(grpcServer, &synchronizationService{})
	go grpcServer.Serve(lis)
	logrus.Infof("Server started. Listening on: %+v", lis.Addr())
}
