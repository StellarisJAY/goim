package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"gorm.io/gorm"
	"time"
)

type GroupServiceImpl struct {
	idGenerator *snowflake.Snowflake
}

func NewGroupServiceImpl() *GroupServiceImpl {
	return &GroupServiceImpl{idGenerator: snowflake.NewSnowflake(config.Config.MachineID)}
}

func (g *GroupServiceImpl) CreateGroup(ctx context.Context, request *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	// 生成群聊ID
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
	// 添加群聊记录
	err = dao.InsertGroup(&model.Group{
		ID:           groupId,
		Name:         request.Name,
		CreateTime:   time.Now().UnixMilli(),
		Description:  request.Description,
		OwnerID:      request.OwnerID,
		OwnerAccount: userInfo.Account,
	})
	if err != nil {
		resp.Code, resp.Message = pb.Error, err.Error()
		return resp, nil
	}
	member := &model.GroupMember{
		GroupID:  groupId,
		UserID:   request.OwnerID,
		JoinTime: time.Now().UnixMilli(),
		Status:   model.MemberStatusNormal,
		Role:     model.MemberRoleOwner,
	}
	// 添加群成员记录
	err = dao.AddGroupMember(member)
	if err != nil {
		resp.Code, resp.Message = pb.Error, err.Error()
		return resp, nil
	}
	resp.GroupID = groupId
	return resp, nil
}

func (g *GroupServiceImpl) GetGroupInfo(ctx context.Context, request *pb.GetGroupInfoRequest) (*pb.GetGroupInfoResponse, error) {
	groupInfo, err := dao.FindGroupInfo(request.GroupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetGroupInfoResponse{Code: pb.NotFound, Message: "group not exists"}, nil
		}
		return &pb.GetGroupInfoResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.GetGroupInfoResponse{
		Code:  pb.Success,
		Group: GroupModelToDTO(groupInfo),
	}, nil
}

func (g *GroupServiceImpl) ListGroupMembers(ctx context.Context, request *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	// 先查看group是否存在
	_, err := dao.FindGroupInfo(request.GroupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ListMembersResponse{Code: pb.NotFound, Message: "group doesn't exist"}, nil
		}
		return &pb.ListMembersResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 查询群成员信息
	members, err := dao.ListGroupMembers(request.GroupID)
	if err != nil {
		return &pb.ListMembersResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	groupMembers := make([]*pb.GroupMember, len(members))
	for i, member := range members {
		groupMembers[i] = &pb.GroupMember{
			UserID:   member.UserID,
			Account:  member.Account,
			NickName: member.NickName,
			JoinTime: member.JoinTime,
			Status:   pb.GroupMemberStatus(member.Status),
		}
	}
	return &pb.ListMembersResponse{
		Code:    pb.Success,
		Total:   int32(len(members)),
		Members: groupMembers,
	}, nil
}

// InviteUser 邀请用户进群
// 先检查用户是否已经进群
// 然后检查inviter是否有邀请权限
// 最后在MongoDB中保留邀请记录
func (g *GroupServiceImpl) InviteUser(ctx context.Context, request *pb.InviteUserRequest) (*pb.InviteUserResponse, error) {
	// 查看用户是否已经进群
	member, err := dao.FindGroupMember(request.GroupID, request.UserID)
	if err == nil && member.UserID != 0 {
		return &pb.InviteUserResponse{Code: pb.InvalidOperation, Message: "member already in group chat"}, nil
	}
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return &pb.InviteUserResponse{Code: pb.Error, Message: err.Error()}, nil
		}
	}
	// 检查邀请者权限
	inviter := dao.FindGroupMemberFull(request.GroupID, request.Inviter)
	if inviter == nil || inviter.Role == model.MemberRoleNormal {
		return &pb.InviteUserResponse{Code: pb.Error, Message: "operation not allowed"}, nil
	}
	noteService, err := getNotificationService()
	if err != nil {
		return &pb.InviteUserResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 添加进群邀请
	addNoteResp, err := noteService.AddNotification(ctx, &pb.AddNotificationRequest{Notification: &pb.Notification{
		Id:          notificationId.NextID(),
		Receiver:    request.UserID,
		TriggerUser: request.GroupID,
		Message:     "",
		Read:        false,
		Type:        int32(model.NotificationGroupInvitation),
		Timestamp:   time.Now().UnixMilli(),
	}})
	if err != nil {
		return &pb.InviteUserResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.InviteUserResponse{Code: addNoteResp.Code, Message: addNoteResp.Message}, nil
}

// AcceptInvitation 接收邀请，进入群聊
func (g *GroupServiceImpl) AcceptInvitation(ctx context.Context, request *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	noteService, err := getNotificationService()
	if err != nil {
		return &pb.AcceptInvitationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 获取通知
	getResp, err := noteService.GetNotification(ctx, &pb.GetNotificationRequest{Id: request.InvitationID})
	if err != nil {
		return &pb.AcceptInvitationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	if getResp.Code != pb.Success {
		return &pb.AcceptInvitationResponse{Code: getResp.Code, Message: getResp.Message}, nil
	}
	// 检查通知是否属于用户 和 通知类型是否正确
	if request.UserID != getResp.Notification.Receiver {
		return &pb.AcceptInvitationResponse{Code: pb.AccessDenied, Message: "userID and notification receiver mismatch"}, nil
	}
	if getResp.Notification.Type != int32(model.NotificationGroupInvitation) {
		return &pb.AcceptInvitationResponse{Code: pb.Error, Message: "invalid group invitation"}, nil
	}
	// 删除通知
	rmResp, err := noteService.RemoveNotification(ctx, &pb.RemoveNotificationRequest{Id: request.InvitationID})
	if err != nil {
		return &pb.AcceptInvitationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	if rmResp.Code != pb.Success {
		return &pb.AcceptInvitationResponse{Code: rmResp.Code, Message: rmResp.Message}, nil
	}

	member := &model.GroupMember{
		GroupID:  getResp.Notification.TriggerUser,
		UserID:   getResp.Notification.Receiver,
		JoinTime: time.Now().UnixMilli(),
		Status:   model.MemberStatusNormal,
		Role:     model.MemberRoleNormal,
	}
	// 添加到群成员列表
	err = dao.AddGroupMember(member)
	if err != nil {
		return &pb.AcceptInvitationResponse{
			Code:    pb.Error,
			Message: err.Error(),
		}, nil
	}
	return &pb.AcceptInvitationResponse{
		Code: pb.Success,
	}, nil
}

// ListGroups 查询用户已经加入的群聊信息列表
func (g *GroupServiceImpl) ListGroups(ctx context.Context, request *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	groupIds, err := dao.ListUserJoinedGroupIds(request.UserID)
	if err != nil {
		return &pb.ListGroupsResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	groups, err := dao.ListGroupInfos(groupIds)
	if err != nil {
		return &pb.ListGroupsResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.ListGroupsResponse{
		Code:    pb.Success,
		Message: "",
		Groups:  GroupsToDTO(groups),
	}, nil
}

func GroupsToDTO(groups []*model.Group) []*pb.GroupInfo {
	infos := make([]*pb.GroupInfo, len(groups))
	for i, g := range groups {
		infos[i] = GroupModelToDTO(g)
	}
	return infos
}

func GroupModelToDTO(group *model.Group) *pb.GroupInfo {
	return &pb.GroupInfo{
		GroupID:      group.ID,
		Name:         group.Name,
		Description:  group.Description,
		OwnerID:      group.OwnerID,
		OwnerAccount: group.OwnerAccount,
		CreateTime:   group.CreateTime,
	}
}
