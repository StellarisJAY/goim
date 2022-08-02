package handler

import (
	_context "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

// CreateGroupHandler 创建群聊处理器
var CreateGroupHandler = func(ctx context.Context) {
	userID := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userID)
	request := &http.CreateGroupRequest{}
	err := ctx.ReadJSON(request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getGroupService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	resp, err := service.CreateGroup(_context.TODO(), &pb.CreateGroupRequest{
		Name:        request.Name,
		OwnerID:     uid,
		Description: request.Description,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	_, _ = ctx.JSON(&http.CreateGroupResponse{
		BaseResponse: http.BaseResponse{
			Code:    resp.Code,
			Message: resp.Message,
		},
		GroupID: resp.GroupID,
	})
}

// GroupInfoHandler 群信息处理器
var GroupInfoHandler = func(ctx context.Context) {
	groupID, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getGroupService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	resp, err := service.GetGroupInfo(_context.TODO(), &pb.GetGroupInfoRequest{GroupID: groupID})
	if err != nil {
		handleError(ctx, err)
		return
	}
	response := &http.GroupInfoResponse{BaseResponse: http.BaseResponse{Code: resp.Code, Message: resp.Message}}
	if resp.Code == pb.Success {
		response.Group = resp.Group
	}
	_, _ = ctx.JSON(response)
}

// GroupMemberHandler 列出群成员处理器
var GroupMemberHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	groupID, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getGroupService()
	if err != nil {
		panic(err)
	}
	resp, err := service.ListGroupMembers(_context.TODO(), &pb.ListMembersRequest{
		GroupID:  groupID,
		PageSize: 0,
		Page:     0,
	})
	if err != nil {
		panic(err)
	}
	response := &http.ListGroupMemberResponse{BaseResponse: http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	}}
	if resp.Code == pb.Success {
		response.Members = resp.Members
	}
	_, _ = ctx.JSON(response)
}

// InviteUserHandler 邀请用户进群处理器
var InviteUserHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	inviteUserID, err := ctx.Params().GetInt64("uid")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	groupID, err := ctx.Params().GetInt64("gid")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	inviterID, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	service, err := getGroupService()
	if err != nil {
		panic(err)
	}
	resp, err := service.InviteUser(_context.TODO(), &pb.InviteUserRequest{
		UserID:  inviteUserID,
		GroupID: groupID,
		Inviter: inviterID,
	})
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(http.BaseResponse{Code: resp.Code, Message: resp.Message})
}

// ListInvitationsHandler 列出用户被邀请记录处理器
var ListInvitationsHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	userID, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	service, err := getGroupService()
	if err != nil {
		panic(err)
	}
	resp, err := service.ListGroupInvitations(_context.TODO(), &pb.ListInvitationRequest{UserID: userID})
	response := &http.ListInvitationsResponse{BaseResponse: http.BaseResponse{}}
	response.Code = resp.Code
	response.Message = resp.Message
	if resp.Code == pb.Success {
		response.Invitations = resp.Invitations
	}
	_, _ = ctx.JSON(response)
}

// AcceptInvitationHandler 接收邀请处理器
var AcceptInvitationHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	invID, err := ctx.Params().GetInt64("invID")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	service, err := getGroupService()
	if err != nil {
		panic(err)
	}
	resp, err := service.AcceptInvitation(_context.TODO(), &pb.AcceptInvitationRequest{
		InvitationID: invID,
		UserID:       uid,
	})
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(&http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	})
}

func getGroupService() (pb.GroupClient, error) {
	conn, err := naming.GetClientConn("user_group")
	if err != nil {
		return nil, err
	}
	return pb.NewGroupClient(conn), nil
}
