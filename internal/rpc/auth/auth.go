package auth

import (
	"github.com/stellarisJAY/goim/internal/rpc/auth/service"
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
	// 注册授权服务
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "auth",
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterAuthServer(server, &service.AuthServiceImpl{})
	log.Info("auth service registered, time used: %d ms", time.Now().Sub(startTime).Milliseconds())
}

func Start() error {
	listen, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		return err
	}
	log.Info("auth service started")
	return server.Serve(listen)
}
