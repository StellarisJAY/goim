package group

import (
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/internal/rpc/group/service"
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
	tracer, closer := trace.NewTracer(pb.GroupServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	server = grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
	startTime := time.Now()
	pb.RegisterGroupServer(server, service.NewGroupServiceImpl())
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: pb.GroupServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("group service registered", zap.Int64("time used(ms)", time.Now().Sub(startTime).Milliseconds()))
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
