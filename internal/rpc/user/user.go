package user

import (
	"github.com/stellarisJAY/goim/internal/rpc/user/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/grpc"
	"net"
	"time"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
	startTime := time.Now()
	pb.RegisterUserServer(server, &service.UserServiceImpl{})
	pb.RegisterAuthServer(server, &service.AuthServiceImpl{})
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: pb.UserServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("user service registered, time used: %d ms", time.Now().Sub(startTime).Milliseconds())
}

func Start() {
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
