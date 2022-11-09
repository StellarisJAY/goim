package user

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/internal/rpc/user/service"
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
	tracer, closer := trace.NewTracer(pb.UserServiceName)
	defer closer.Close()
	server = grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
	startTime := time.Now()
	pb.RegisterUserServer(server, &service.UserServiceImpl{})
	pb.RegisterAuthServer(server, &service.AuthServiceImpl{})
	err := naming.RegisterService(naming.ServiceRegistration{
		ID:          "",
		ServiceName: pb.UserServiceName,
		Address:     config.Config.RpcServer.Address,
	})
	if err != nil {
		panic(err)
	}
	log.Info("user service registered",
		zap.Int64("time used(ms)", time.Now().Sub(startTime).Milliseconds()))
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(config.Config.Metrics.PromHttpAddr, nil)
	}()
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
