package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

type UserServiceImpl struct {
}

func (u *UserServiceImpl) FindUserByID(ctx context.Context, request *pb.FindUserByIdRequest) (*pb.FindUserByIdResponse, error) {
	userInfo, err := dao.FindUserInfo(request.Id)
	if err != nil {
		return &pb.FindUserByIdResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.FindUserByIdResponse{Code: pb.Success, User: &pb.UserInfo{
		Id:           userInfo.ID,
		Account:      userInfo.Account,
		NickName:     userInfo.NickName,
		RegisterTime: userInfo.RegisterTime,
	}}, nil
}

func (u *UserServiceImpl) UpdateUserInfo(ctx context.Context, request *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	err := dao.UpdateUserInfo(&model.UserInfo{
		ID:       request.Id,
		NickName: request.NickName,
	})
	if err != nil {
		return &pb.UpdateUserInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.UpdateUserInfoResponse{Code: pb.Success}, nil
}
