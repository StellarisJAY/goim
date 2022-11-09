package friend

import (
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/internal/rpc/friend/service"
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
	tracer, closer := trace.NewTracer(pb.FriendServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	server = grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
	startTime := time.Now()
	pb.RegisterFriendServer(server, service.NewFriendServiceImpl())
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: pb.FriendServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("friend service registered", zap.Int64("time used(ms)", time.Now().Sub(startTime).Milliseconds()))
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(config.Config.Metrics.PromHttpAddr, nil)
	}()
	go func() {
		service.AsyncPushSendNotification()
	}()
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
