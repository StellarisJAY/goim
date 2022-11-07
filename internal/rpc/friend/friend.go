package friend

import (
	"github.com/stellarisJAY/goim/internal/rpc/friend/service"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
)

var server *grpc.Server

func Init() {
	server = grpc.NewServer()
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
}

func Start() {
	listener, err := net.Listen("tcp", config.Config.RpcServer.Address)
	if err != nil {
		panic(err)
	}
	go func() {
		service.AsyncPushSendNotification()
	}()
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
