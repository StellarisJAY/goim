package service

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/cache"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"time"
)

type FriendServiceImpl struct {
	applicationId *snowflake.Snowflake
}

func NewFriendServiceImpl() *FriendServiceImpl {
	return &FriendServiceImpl{applicationId: snowflake.NewSnowflake(config.Config.MachineID)}
}

func (f *FriendServiceImpl) AddFriend(ctx context.Context, request *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	// 查询目标用户是否存在
	_, err := dao.FindUserInfo(request.TargetUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.AddFriendResponse{Code: pb.NotFound, Message: "target user not found"}, nil
		}
		return &pb.AddFriendResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 查询好友关系是否已经存在
	friendship, err := dao.CheckFriendship(request.UserID, request.TargetUser)
	if err != nil {
		return &pb.AddFriendResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	if friendship {
		return &pb.AddFriendResponse{
			Code:    pb.InvalidOperation,
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
		if mongo.IsDuplicateKeyError(err) {
			return &pb.AddFriendResponse{Code: pb.InvalidOperation, Message: "duplicate application"}, nil
		}
		return &pb.AddFriendResponse{Code: pb.Error, Message: "can't create add friend application"}, nil
	}
	return &pb.AddFriendResponse{
		Code: pb.Success,
	}, nil
}

func (f *FriendServiceImpl) ListAddFriendRequests(ctx context.Context, request *pb.ListAddFriendRequest) (*pb.ListAddFriendResponse, error) {
	applications, err := dao.ListAddFriendRequests(request.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.ListAddFriendResponse{Code: pb.NotFound, Message: "no application found"}, nil
		}
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
		if err == mongo.ErrNoDocuments {
			return &pb.AcceptFriendResponse{Code: pb.NotFound, Message: "no such application found"}, nil
		}
		return &pb.AcceptFriendResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	acceptTime := time.Now().UnixMilli()
	fs1 := &model.Friend{
		OwnerID:    application.Requester,
		FriendID:   application.Target,
		AcceptTime: acceptTime,
	}
	fs2 := &model.Friend{
		OwnerID:    application.Target,
		FriendID:   application.Requester,
		AcceptTime: acceptTime,
	}
	// MySQL记录好友关系
	err = dao.InsertFriendship(fs1, fs2)
	// 删除缓存
	_ = cache.Delete(fmt.Sprintf("%s%d", dao.KeyFriendIDList, application.Target))
	if err != nil {
		return &pb.AcceptFriendResponse{
			Code:    pb.Error,
			Message: err.Error(),
		}, nil
	}
	return &pb.AcceptFriendResponse{Code: pb.Success}, nil
}

func (f *FriendServiceImpl) ListFriends(ctx context.Context, request *pb.FriendListRequest) (*pb.FriendListResponse, error) {
	friendIDs, err := dao.ListFriendIDs(request.UserID)
	if err != nil {
		return &pb.FriendListResponse{Code: pb.Error, Message: err.Error()}, nil
	}

	infos := make([]*pb.FriendInfo, 0, len(friendIDs))
	for _, friend := range friendIDs {
		userInfo, err := dao.FindUserInfo(friend)
		if err != nil {
			continue
		}
		infos = append(infos, &pb.FriendInfo{
			UserID:       userInfo.ID,
			Account:      userInfo.Account,
			NickName:     userInfo.NickName,
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
	friendship, err := dao.CheckFriendship(request.UserID, request.FriendID)
	if err != nil {
		return &pb.FriendInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	if !friendship {
		return &pb.FriendInfoResponse{Code: pb.AccessDenied, Message: "target user is not your friend"}, nil
	}
	// 通过好友关系查询好友的个人信息
	friendInfo, err := dao.FindUserInfo(request.FriendID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.FriendInfoResponse{Code: pb.NotFound, Message: "friend information not found"}, nil
		}
		return &pb.FriendInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.FriendInfoResponse{
		Code:    pb.Success,
		Message: "",
		Info: &pb.FriendInfo{
			UserID:       friendInfo.ID,
			Account:      friendInfo.Account,
			NickName:     friendInfo.NickName,
			RegisterTime: friendInfo.RegisterTime,
		},
	}, nil
}
