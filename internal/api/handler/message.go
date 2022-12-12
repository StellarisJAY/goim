package handler

import (
	_context "context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.SyncOfflineGroupMessageRequest{})
}

// SyncOfflineMessageHandler 同步离线消息
var SyncOfflineMessageHandler = func(ctx *context.Context) {
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
	_ = ctx.JSON(resp)
}

var SyncOfflineGroupMessages = func(ctx *context.Context) {
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
		_ = ctx.JSON(response)
	}
}

var SyncGroupMessages = func(ctx *context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	userID := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userID)
	groupID, err := ctx.Params().GetInt64("groupID")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid seq number")
		return
	}
	lastSeq, err := ctx.Params().GetInt64("seq")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid seq number")
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	resp, err := service.SyncGroupMessages(_context.TODO(), &pb.SyncGroupMessageRequest{
		GroupID: groupID,
		Seq:     lastSeq,
		UserID:  uid,
	})
	if err != nil {
		panic(err)
	} else {
		response := &http.SyncGroupMessageResponse{
			BaseResponse: http.BaseResponse{Code: resp.Code, Message: resp.Message},
			Messages:     resp.Msgs,
		}
		_ = ctx.JSON(response)
	}
}

var SyncLatestGroupMessages = func(ctx *context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	userID := ctx.Params().Get("userID")
	uid, _ := stringutil.HexStringToInt64(userID)
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
	resp, err := service.SyncLatestGroupMessages(_context.TODO(), &pb.SyncLatestGroupMessageRequest{
		GroupID: request.GroupID,
		Limit:   request.Limit,
		UserID:  uid,
	})
	if err != nil {
		panic(err)
	}
	response := &http.SyncGroupLatestMessagesResponse{
		BaseResponse: http.BaseResponse{Code: resp.Code, Message: resp.Message},
		FirstSeq:     resp.FirstSeq,
		LastSeq:      resp.LastSeq,
		Msgs:         resp.Msgs,
	}
	_ = ctx.JSON(response)
}

var ListNotifications = func(ctx *context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	notReadOption, err := ctx.Params().GetBool("notRead")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("wrong read option")
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	response, err := service.ListNotifications(_context.TODO(), &pb.ListNotificationRequest{UserID: uid, NotRead: notReadOption})
	if err != nil {
		panic(err)
	}
	resp := &http.ListNotificationRequest{
		BaseResponse: http.BaseResponse{
			Code:    response.Code,
			Message: response.Message,
		},
		Notifications: convertNotificationStruct(response.Notifications),
	}
	_ = ctx.JSON(resp)
}

var MarkNotificationReadHandler = func(ctx *context.Context) {
	defer func() {
		if err, ok := recover().(error); err != nil && ok {
			handleError(ctx, err)
		}
	}()
	uid, _ := stringutil.HexStringToInt64(ctx.Params().Get("userID"))
	notificationID, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid notification id")
		return
	}
	service, err := getMessageService()
	if err != nil {
		panic(err)
	}
	response, err := service.MarkNotificationRead(_context.TODO(), &pb.MarkNotificationReadRequest{NotificationID: notificationID, UserID: uid})
	if err != nil {
		panic(err)
	}
	_ = ctx.JSON(&http.BaseResponse{Code: response.Code, Message: response.Message})
}

func convertNotificationStruct(pbNotifications []*pb.Notification) []*http.Notification {
	notifications := make([]*http.Notification, len(pbNotifications))
	for i, notification := range pbNotifications {
		notifications[i] = &http.Notification{
			Id:          notification.Id,
			Receiver:    notification.Receiver,
			TriggerUser: notification.TriggerUser,
			Type:        byte(notification.Type),
			Message:     notification.Message,
			Read:        notification.Read,
			Timestamp:   notification.Timestamp,
		}
	}
	return notifications
}
