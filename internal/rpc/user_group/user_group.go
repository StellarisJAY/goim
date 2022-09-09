package user_group

import (
	"github.com/stellarisJAY/goim/internal/rpc/user_group/service"
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
	pb.RegisterGroupServer(server, service.NewGroupServiceImpl())
	pb.RegisterFriendServer(server, service.NewFriendServiceImpl())
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: "user_group",
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("user-group service registered, time used: %d ms", time.Now().Sub(startTime).Milliseconds())
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
