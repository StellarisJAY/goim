package msg

import (
	"github.com/stellarisJAY/goim/internal/rpc/msg/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
	pb.RegisterMessageServer(server, service.NewMessageServiceImpl())
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
		ServiceName: "message",
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
}

func Start() {
	listen, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		log.Println(err)
	}
}
