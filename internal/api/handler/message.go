package handler

import (
	_context "context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.SyncOfflineGroupMessageRequest{})
}

// SyncOfflineMessageHandler 同步离线消息
var SyncOfflineMessageHandler = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	userID := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userID)
	// 客户端保存的最大的消息序号
	seq, err := ctx.Params().GetInt64("seq")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid sequence number: " + err.Error())
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	// RPC 查询客户端最大序号之后的所有消息
	response, err := service.SyncOfflineMessages(_context.TODO(), &pb.SyncMsgRequest{
		LastSeq: seq,
		UserID:  uid,
	})
	if err != nil {
		panic(err)
	}
	resp := new(http.SyncOfflineMessageResponse)
	resp.BaseResponse = http.BaseResponse{}
	if response.Code == pb.Success {
		resp.BaseResponse.Code = response.Code
		resp.LastSeq = response.LastSeq
		resp.InitSeq = response.InitSeq
		resp.Messages = response.Messages
	} else {
		resp.BaseResponse.Code = response.Code
		resp.BaseResponse.Message = response.Message
	}
	_, _ = ctx.JSON(resp)
}

var SyncOfflineGroupMessages = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	request := http.SyncOfflineGroupMessageRequest{}
	err := ctx.ReadJSON(&request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid request")
		return
	}
	if err := validate.Struct(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(fmt.Sprintf("request validate failed, error: %v", err))
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	resp, err := service.SyncOfflineGroupMessages(_context.TODO(), &pb.SyncGroupMsgRequest{
		Groups:     request.Groups,
		Timestamps: request.Timestamps,
	})
	if err != nil {
		panic(err)
	} else {
		response := &http.SyncOfflineGroupMessageResponse{
			BaseResponse:  http.BaseResponse{Code: resp.Code, Message: resp.Message},
			GroupMessages: resp.GroupMessages,
		}
		_, _ = ctx.JSON(response)
	}
}

var SyncLatestGroupMessages = func(ctx context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	request := new(http.SyncGroupLatestMessagesRequest)
	err := ctx.ReadJSON(request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("bad request")
		return
	}
	if err = validate.Struct(request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("bad request")
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	resp, err := service.SyncGroupLatestMessages(_context.TODO(), &pb.SyncGroupLatestMessagesRequest{
		GroupID:       request.GroupID,
		Limit:         request.Limit,
		LastTimestamp: request.LastTimestamp,
	})
	if err != nil {
		panic(err)
	}
	response := &http.SyncGroupLatestMessagesResponse{
		BaseResponse: http.BaseResponse{Code: resp.Code, Message: resp.Message},
		Msgs:         resp.Msgs,
	}
	_, _ = ctx.JSON(response)
}

func getMessageService() (pb.MessageClient, error) {
	conn, err := naming.GetClientConn("message")
	if err != nil {
		return nil, err
	}
	client := pb.NewMessageClient(conn)
	return client, nil
}
