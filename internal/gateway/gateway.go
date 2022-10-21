package gateway

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/mq/kafka"
	"github.com/stellarisJAY/goim/pkg/mq/nsq"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/grpc"
	"net"
	"strings"
	"time"
)

type Server struct {
	wsServer    *websocket.Server
	grpcServer  *grpc.Server
	kafkaCG     *kafka.ConsumerGroup
	nsqConsumer *nsq.Consumer
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
	useMQ := strings.ToLower(config.Config.MessageQueue)
	group := config.Config.Gateway.ConsumerGroup
	if group == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			panic(fmt.Errorf("auto generate gateway uuid failed error %w", err))
		}
		group = id.String()
	}
	switch useMQ {
	case "kafka":
		s.kafkaCG = kafka.NewConsumerGroup(group, config.Config.Kafka.Addrs, []string{pb.MessagePushTopic})
		s.kafkaCG.Start(context.TODO(), s.HandleKafka)
	case "nsq":
		//todo gateway 单独的channel
		s.nsqConsumer = nsq.NewConsumer(pb.MessagePushTopic, group, s.HandleNSQ)
		s.nsqConsumer.Connect()
	}
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
