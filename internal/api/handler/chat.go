package handler

import (
	_context "context"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/copier"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
)

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, &http.SendMessageRequest{})
}

var SendMessageHandler = func(ctx context.Context) {
	userID := ctx.Params().Get("userID")
	req := new(http.SendMessageRequest)
	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	message := new(pb.BaseMsg)
	if err := copier.CopyStructFields(message, req); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	uid, _ := stringutil.HexStringToInt64(userID)
	message.From = uid
	conn, err := naming.GetClientConn("chat")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	client := pb.NewChatClient(conn)
	response, err := client.SendMessage(_context.Background(), &pb.SendMsgRequest{Msg: message})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(&http.SendMessageResponse{
		BaseResponse: http.BaseResponse{Code: response.Code, Message: response.Message},
		MessageID:    response.MessageId,
		Timestamp:    response.Timestamp,
	})
}
