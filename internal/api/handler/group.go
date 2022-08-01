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

var GroupMemberHandler = func(ctx context.Context) {

}

var InviteUserHandler = func(ctx context.Context) {

}

func getGroupService() (pb.GroupClient, error) {
	conn, err := naming.GetClientConn("user_group")
	if err != nil {
		return nil, err
	}
	return pb.NewGroupClient(conn), nil
}
