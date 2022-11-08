package handler

import (
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

func getUserService() (pb.UserClient, error) {
	conn, err := naming.GetClientConn(pb.UserServiceName)
	if err != nil {
		return nil, err
	}
	return pb.NewUserClient(conn), nil
}

func getMessageService() (pb.MessageClient, error) {
	conn, err := naming.GetClientConn(pb.MessageServiceName)
	if err != nil {
		return nil, err
	}
	client := pb.NewMessageClient(conn)
	return client, nil
}

func getGroupService() (pb.GroupClient, error) {
	conn, err := naming.GetClientConn(pb.GroupServiceName)
	if err != nil {
		return nil, err
	}
	return pb.NewGroupClient(conn), nil
}

func getFriendService() (pb.FriendClient, error) {
	conn, err := naming.GetClientConn(pb.FriendServiceName)
	if err != nil {
		return nil, err
	}
	return pb.NewFriendClient(conn), nil
}

func GetAuthService() (pb.AuthClient, error) {
	// 从服务发现获取 RPC 客户端连接
	conn, err := naming.GetClientConn(pb.UserServiceName)
	if err != nil {
		return nil, err
	}
	// RPC调用用户注册服务
	return pb.NewAuthClient(conn), nil
}
