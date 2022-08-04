package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"time"
)

type FriendServiceImpl struct {
	applicationId *snowflake.Snowflake
}

func NewFriendServiceImpl() *FriendServiceImpl {
	return &FriendServiceImpl{applicationId: snowflake.NewSnowflake(config.Config.MachineID)}
}

func (f *FriendServiceImpl) AddFriend(ctx context.Context, request *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	_, err := dao.FindUserInfo(request.TargetUser)
	if err != nil {
		return &pb.AddFriendResponse{Code: pb.Error, Message: "target user not available"}, nil
	}
	friendship, err := dao.GetFriendInfo(request.UserID, request.TargetUser)
	if friendship != nil {
		return &pb.AddFriendResponse{
			Code:    pb.Error,
			Message: "already established friendship",
		}, nil
	}
	err = dao.AddFriendRequest(&model.AddFriendRequest{
		Requester: request.UserID,
		Target:    request.TargetUser,
		Timestamp: time.Now().UnixMilli(),
		Message:   request.Message,
	})
	if err != nil {
		return &pb.AddFriendResponse{Code: pb.Error, Message: "can't create add friend application"}, nil
	}
	return &pb.AddFriendResponse{
		Code: pb.Success,
	}, nil
}

func (f *FriendServiceImpl) ListAddFriendRequests(ctx context.Context, request *pb.ListAddFriendRequest) (*pb.ListAddFriendResponse, error) {
	applications, err := dao.ListAddFriendRequests(request.UserID)
	if err != nil {
		return &pb.ListAddFriendResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	friendApplications := make([]*pb.FriendApplication, 0, len(applications))
	for _, app := range applications {
		requester, err := dao.FindUserInfo(app.Requester)
		if err != nil {
			continue
		}
		friendApplications = append(friendApplications, &pb.FriendApplication{
			UserID:    requester.ID,
			Account:   requester.Account,
			NickName:  requester.NickName,
			Timestamp: app.Timestamp,
			Message:   app.Message,
		})
	}
	return &pb.ListAddFriendResponse{
		Code:         pb.Success,
		Applications: friendApplications,
	}, nil
}

func (f *FriendServiceImpl) AcceptFriend(ctx context.Context, request *pb.AcceptFriendRequest) (*pb.AcceptFriendResponse, error) {
	application, err := dao.GetAndDeleteFriendRequest(request.TargetID, request.UserID)
	if err != nil {
		return &pb.AcceptFriendResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	acceptTime := time.Now().UnixMilli()
	err = dao.InsertFriendship(&model.Friend{
		OwnerID:    application.Requester,
		FriendID:   application.Target,
		AcceptTime: acceptTime,
	}, &model.Friend{
		OwnerID:    application.Target,
		FriendID:   application.Requester,
		AcceptTime: acceptTime,
	})
	if err != nil {
		return &pb.AcceptFriendResponse{
			Code:    pb.Error,
			Message: err.Error(),
		}, nil
	}
	return &pb.AcceptFriendResponse{Code: pb.Success}, nil
}

func (f *FriendServiceImpl) ListFriends(ctx context.Context, request *pb.FriendListRequest) (*pb.FriendListResponse, error) {
	friends, err := dao.ListFriends(request.UserID)
	if err != nil {
		return &pb.FriendListResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	infos := make([]*pb.FriendInfo, 0, len(friends))
	for _, friend := range friends {
		userInfo, err := dao.FindUserInfo(friend.FriendID)
		if err != nil {
			continue
		}
		infos = append(infos, &pb.FriendInfo{
			UserID:       userInfo.ID,
			Account:      userInfo.Account,
			NickName:     userInfo.NickName,
			AcceptTime:   friend.AcceptTime,
			RegisterTime: userInfo.RegisterTime,
		})
	}
	return &pb.FriendListResponse{
		Code:    pb.Success,
		Message: "",
		Friends: infos,
	}, nil
}

func (f *FriendServiceImpl) GetFriendInfo(ctx context.Context, request *pb.FriendInfoRequest) (*pb.FriendInfoResponse, error) {
	// 查询好友关系是否存在
	friend, err := dao.GetFriendInfo(request.UserID, request.FriendID)
	if err != nil {
		return &pb.FriendInfoResponse{
			Code:    pb.Error,
			Message: err.Error(),
		}, nil
	}
	// 通过好友关系查询好友的个人信息
	friendInfo, err := dao.FindUserInfo(friend.FriendID)
	if err != nil {
		return &pb.FriendInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.FriendInfoResponse{
		Code:    pb.Success,
		Message: "",
		Info: &pb.FriendInfo{
			UserID:       friendInfo.ID,
			Account:      friendInfo.Account,
			NickName:     friendInfo.NickName,
			AcceptTime:   friend.AcceptTime,
			RegisterTime: friendInfo.RegisterTime,
		},
	}, nil
}
