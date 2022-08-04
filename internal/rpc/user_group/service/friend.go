package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

type FriendServiceImpl struct {
}

func (f *FriendServiceImpl) AddFriend(ctx context.Context, request *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FriendServiceImpl) ListAddFriendRequests(ctx context.Context, request *pb.ListAddFriendRequest) (*pb.ListAddFriendResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FriendServiceImpl) AcceptFriend(ctx context.Context, request *pb.AcceptFriendRequest) (*pb.AcceptFriendResponse, error) {
	//TODO implement me
	panic("implement me")
}
