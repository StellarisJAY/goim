package handler

import (
	_context "context"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
	"github.com/stellarisJAY/goim/pkg/trace"
)

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.SendMessageRequest{})
}

var SendMessageHandler = func(ctx context.Context) {
	userID := ctx.Params().Get("userID")
	deviceID := ctx.Params().Get("deviceID")
	req := new(http.SendMessageRequest)
	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	message := &pb.BaseMsg{
		To:       req.To,
		Content:  req.Content,
		Flag:     pb.MessageFlag(req.Flag),
		DeviceId: deviceID,
	}
	uid, _ := stringutil.HexStringToInt64(userID)
	message.From = uid
	tracer, closer := trace.NewTracer("api-chat-handler")
	defer closer.Close()
	service, err := getMessageService(tracer)
	if err != nil {
		handleError(ctx, err)
		return
	}
	response, err := service.SendMessage(_context.Background(), &pb.SendMsgRequest{Msg: message})
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(&http.SendMessageResponse{
		BaseResponse: http.BaseResponse{Code: response.Code, Message: response.Message},
		MessageID:    response.MessageId,
		Timestamp:    response.Timestamp,
	})
}
