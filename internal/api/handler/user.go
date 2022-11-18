package handler

import (
	_context "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

// FindUserHandler 通过ID获取用户基本信息
var FindUserHandler = func(ctx *context.Context) {
	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	service, err := getUserService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	findUserResponse, err := service.FindUserByID(_context.TODO(), &pb.FindUserByIdRequest{Id: id})
	if err != nil {
		handleError(ctx, err)
		return
	}
	response := &http.FindUserResponse{}
	response.Code = findUserResponse.Code
	response.Message = findUserResponse.Message
	response.UserInfo = findUserResponse.User
	_ = ctx.JSON(response)
}

var GetSelfInfoHandler = func(ctx *context.Context) {
	userId := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userId)
	service, err := getUserService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	findUserResponse, err := service.FindUserByID(_context.TODO(), &pb.FindUserByIdRequest{Id: uid})
	if err != nil {
		handleError(ctx, err)
		return
	}
	response := &http.FindUserResponse{}
	response.Code = findUserResponse.Code
	response.Message = findUserResponse.Message
	response.UserInfo = findUserResponse.User
	_ = ctx.JSON(response)
}

// UpdateUserHandler 更新用户信息
var UpdateUserHandler = func(ctx *context.Context) {
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
		handleError(ctx, err)
		return
	}
	resp, err := service.UpdateUserInfo(_context.TODO(), &pb.UpdateUserInfoRequest{Id: uid, NickName: req.NickName})
	if err != nil {
		handleError(ctx, err)
		return
	}
	_ = ctx.JSON(&http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	})
}
