package user_group

import (
	"github.com/stellarisJAY/goim/internal/rpc/user_group/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/grpc"
	"net"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
	pb.RegisterUserServer(server, &service.UserServiceImpl{})
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: "user_group",
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
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
