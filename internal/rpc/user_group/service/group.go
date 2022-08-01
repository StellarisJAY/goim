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

type GroupServiceImpl struct {
	idGenerator *snowflake.Snowflake
}

func NewGroupServiceImpl() *GroupServiceImpl {
	return &GroupServiceImpl{idGenerator: snowflake.NewSnowflake(config.Config.MachineID)}
}

func (g *GroupServiceImpl) CreateGroup(ctx context.Context, request *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	groupId := g.idGenerator.NextID()
	resp := &pb.CreateGroupResponse{}
	userInfo, err := dao.FindUserInfo(request.OwnerID)
	if err != nil {
		resp.Code = pb.Error
		resp.Message = err.Error()
		return resp, nil
	}
	if userInfo == nil {
		resp.Code = pb.NotFound
		resp.Message = "owner info not found"
		return resp, nil
	}
	err = dao.InsertGroup(&model.Group{
		ID:           groupId,
		Name:         request.Name,
		CreateTime:   time.Now().UnixMilli(),
		Description:  request.Description,
		OwnerID:      request.OwnerID,
		OwnerAccount: userInfo.Account,
	})
	if err != nil {
		resp.Code = pb.Error
		resp.Message = err.Error()
	} else {
		resp.GroupID = groupId
	}
	return resp, nil
}

func (g *GroupServiceImpl) GetGroupInfo(ctx context.Context, request *pb.GetGroupInfoRequest) (*pb.GetGroupInfoResponse, error) {
	groupInfo, err := dao.FindGroupInfo(request.GroupID)
	if err != nil {
		return &pb.GetGroupInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.GetGroupInfoResponse{
		Code: pb.Success,
		Group: &pb.GroupInfo{
			GroupID:      groupInfo.ID,
			Name:         groupInfo.Name,
			Description:  groupInfo.Description,
			OwnerID:      groupInfo.OwnerID,
			OwnerAccount: groupInfo.OwnerAccount,
			CreateTime:   groupInfo.CreateTime,
		},
	}, nil
}

func (g *GroupServiceImpl) ListGroupMembers(ctx context.Context, request *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	panic("implement me")
}

func (g *GroupServiceImpl) InviteUser(ctx context.Context, request *pb.InviteUserRequest) (*pb.InviteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupServiceImpl) AcceptInvitation(ctx context.Context, request *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}
