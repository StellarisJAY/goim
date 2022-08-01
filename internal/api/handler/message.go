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

// SyncOfflineMessageHandler 同步离线消息
var SyncOfflineMessageHandler = func(ctx context.Context) {
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
		handleError(ctx, err)
		return
	}
	// RPC 查询客户端最大序号之后的所有消息
	response, err := service.SyncOfflineMessages(_context.TODO(), &pb.SyncMsgRequest{
		LastSeq: seq,
		UserID:  uid,
	})
	if err != nil {
		handleError(ctx, err)
		return
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

func getMessageService() (pb.MessageClient, error) {
	conn, err := naming.GetClientConn("message")
	if err != nil {
		return nil, err
	}
	client := pb.NewMessageClient(conn)
	return client, nil
}
