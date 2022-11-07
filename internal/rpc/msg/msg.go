package msg

import (
	"github.com/stellarisJAY/goim/internal/rpc/msg/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
	pb.RegisterMessageServer(server, service.NewMessageServiceImpl())
	startTime := time.Now()
	// 注册聊天服务
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "chat",
		Address:     config.Config.RpcServer.Address,
	})
	// 注册消息查询服务
	if err != nil {
		panic(err)
	}
	err = naming.RegisterService(naming.ServiceRegistration{
		ServiceName: pb.MessageServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("message service registered",
		zap.Int64("time used(ms)", time.Now().Sub(startTime).Milliseconds()))
}

func Start() {
	listen, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
