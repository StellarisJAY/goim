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

// FindUserHandler 通过ID获取用户基本信息
var FindUserHandler = func(ctx context.Context) {
	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	service, err := getUserService()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	findUserResponse, err := service.FindUserByID(_context.TODO(), &pb.FindUserByIdRequest{Id: id})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	response := &http.FindUserResponse{}
	response.Code = findUserResponse.Code
	response.Message = findUserResponse.Message
	response.UserInfo = findUserResponse.User
	_, _ = ctx.JSON(response)
}

// UpdateUserHandler 更新用户信息
var UpdateUserHandler = func(ctx context.Context) {
	userId := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userId)
	req := &http.UpdateUserRequest{}
	err := ctx.ReadJSON(req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getUserService()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	resp, err := service.UpdateUserInfo(_context.TODO(), &pb.UpdateUserInfoRequest{Id: uid, NickName: req.NickName})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	_, _ = ctx.JSON(&http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	})
}

func getUserService() (pb.UserClient, error) {
	conn, err := naming.GetClientConn("user_group")
	if err != nil {
		return nil, err
	}
	return pb.NewUserClient(conn), nil
}
