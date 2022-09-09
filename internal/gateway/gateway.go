package gateway

import (
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	wsServer   *websocket.Server
	grpcServer *grpc.Server
}

func (s *Server) Init() {
	s.grpcServer = grpc.NewServer()
	startTime := time.Now().UnixMilli()
	// 注册网关服务，提供消息下行的RPC接口
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: "gateway",
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterRelayServer(s.grpcServer, s)
	s.wsServer = websocket.NewServer(config.Config.WebsocketServer.Address)
	s.wsServer.Acceptor = &GateAcceptor{}
	log.Info("Gateway service registered, time used: %d ms", time.Now().UnixMilli()-startTime)
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
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
	log.Info("gateway server started")
	return nil
}
