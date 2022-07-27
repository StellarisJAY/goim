package auth

import (
	"github.com/stellarisJAY/goim/internal/rpc/auth/service"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/grpc"
	"net"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
	// 注册授权服务
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "auth",
		Address:     "127.0.0.1:9999",
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterAuthServer(server, &service.AuthServiceImpl{})
}

func Start() error {
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		return err
	}
	return server.Serve(listen)
}
