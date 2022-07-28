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
	chatService := service.NewChatServiceImpl()
	pb.RegisterChatServer(server, chatService)
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "chat",
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
