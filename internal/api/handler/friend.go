package handler

import (
	_context "context"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.AddFriendRequest{})
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.AcceptFriendRequest{})
}

var InsertFriendApplicationHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	request := &http.AddFriendRequest{}
	if err := ctx.ReadJSON(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getFriendService()
	if err != nil {
		panic(err)
	}
	resp, err := service.AddFriend(_context.TODO(), &pb.AddFriendRequest{
		UserID:     uid,
		TargetUser: request.TargetID,
		Message:    request.ValidateMessage,
	})
	if err != nil {
		panic(err)
	}
	response := &http.BaseResponse{}
	if resp.Code == pb.DuplicateKey {
		response.Code = pb.DuplicateKey
		response.Message = "duplicate friend application"
	} else {
		response.Code = resp.Code
		response.Message = resp.Message
	}
	_, _ = ctx.JSON(response)
}

var AcceptFriendHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	request := &http.AcceptFriendRequest{}
	if err := ctx.ReadJSON(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getFriendService()
	if err != nil {
		panic(err)
	}
	resp, err := service.AcceptFriend(_context.TODO(), &pb.AcceptFriendRequest{
		UserID:   uid,
		TargetID: request.TargetID,
	})
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(&http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	})
}

var FriendListHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	service, err := getFriendService()
	if err != nil {
		panic(err)
	}
	resp, err := service.ListFriends(_context.TODO(), &pb.FriendListRequest{UserID: uid})
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(&http.ListFriendsResponse{BaseResponse: http.BaseResponse{
		Code:    resp.Code,
		Message: resp.Message,
	}, Friends: resp.Friends})
}

var FriendInfoHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	friendID, err := ctx.Params().GetInt64("friendID")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := getFriendService()
	if err != nil {
		panic(err)
	}
	resp, err := service.GetFriendInfo(_context.TODO(), &pb.FriendInfoRequest{
		UserID:   uid,
		FriendID: friendID,
	})
	if err != nil {
		panic(err)
	}
	_, _ = ctx.JSON(&http.GetFriendInfoResponse{
		BaseResponse: http.BaseResponse{
			Code:    resp.Code,
			Message: resp.Message,
		},
		Info: resp.Info,
	})
}

func getFriendService() (pb.FriendClient, error) {
	conn, err := naming.GetClientConn(pb.FriendServiceName)
	if err != nil {
		return nil, err
	}
	return pb.NewFriendClient(conn), nil
}
