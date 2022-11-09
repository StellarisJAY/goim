package handler

import (
	"github.com/opentracing/opentracing-go"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

func getUserService(tracer opentracing.Tracer) (pb.UserClient, error) {
	conn, err := naming.GetClientConn(pb.UserServiceName, tracer)
	if err != nil {
		return nil, err
	}
	return pb.NewUserClient(conn), nil
}

func getMessageService(tracer opentracing.Tracer) (pb.MessageClient, error) {
	conn, err := naming.GetClientConn(pb.MessageServiceName, tracer)
	if err != nil {
		return nil, err
	}
	client := pb.NewMessageClient(conn)
	return client, nil
}

func getGroupService(tracer opentracing.Tracer) (pb.GroupClient, error) {
	conn, err := naming.GetClientConn(pb.GroupServiceName, tracer)
	if err != nil {
		return nil, err
	}
	return pb.NewGroupClient(conn), nil
}

func getFriendService(tracer opentracing.Tracer) (pb.FriendClient, error) {
	conn, err := naming.GetClientConn(pb.FriendServiceName, tracer)
	if err != nil {
		return nil, err
	}
	return pb.NewFriendClient(conn), nil
}

func GetAuthService(tracer opentracing.Tracer) (pb.AuthClient, error) {
	// 从服务发现获取 RPC 客户端连接
	conn, err := naming.GetClientConn(pb.UserServiceName, tracer)
	if err != nil {
		return nil, err
	}
	// RPC调用用户注册服务
	return pb.NewAuthClient(conn), nil
}
