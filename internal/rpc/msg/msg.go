package msg

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/internal/rpc/msg/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

var server *grpc.Server

func Init() {

}

func Start() {
	tracer, closer := trace.NewTracer(pb.MessageServiceName)
	defer closer.Close()
	server = grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
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
	listen, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(config.Config.Metrics.PromHttpAddr, nil)
	}()
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
