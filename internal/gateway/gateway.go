package gateway

import (
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	wsServer   *websocket.Server
	grpcServer *grpc.Server
}

func (s *Server) Init() {
	s.grpcServer = grpc.NewServer()
	// 注册网关服务，提供消息下行的RPC接口
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "gateway",
		Address:     "127.0.0.1:9000",
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterRelayServer(s.grpcServer, s)
	s.wsServer = websocket.NewServer(":8000")
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}
	// 启动 RPC 和 Websocket
	go func() {
		_ = s.grpcServer.Serve(listener)
	}()
	go func() {
		_ = s.wsServer.Start()
	}()
	return nil
}
