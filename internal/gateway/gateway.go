package gateway

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/mq/kafka"
	"github.com/stellarisJAY/goim/pkg/mq/nsq"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/trace"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
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
		s.nsqConsumer = nsq.NewConsumer(pb.MessagePushTopic, group, s.HandleNSQ)
		s.nsqConsumer.Connect()
	}
}

func (s *Server) Start() error {
	tracer, closer := trace.NewTracer(pb.GatewayServiceName)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
	startTime := time.Now().UnixMilli()
	// 注册网关服务，提供消息下行的RPC接口
	err := naming.RegisterService(naming.ServiceRegistration{
		ServiceName: pb.GatewayServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterRelayServer(s.grpcServer, s)
	log.Info("Gateway service registered", zap.Int64("time used(ms)", time.Now().UnixMilli()-startTime))

	s.wsServer = websocket.NewServer(config.Config.WebsocketServer.Address)
	s.wsServer.Acceptor = &GateAcceptor{globalTracer: tracer}
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		return err
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(config.Config.Metrics.PromHttpAddr, nil)
	}()
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
