package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/authutil"
	"log"
)

const (
	CtxKeyHandlerError = "handler_error"
)

var TokenVerifier = func(ctx context.Context) {
	var token string
	if token = ctx.GetHeader("Authorization"); token == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.EndRequest()
		return
	}
	userID, deviceID, valid := authutil.ValidateToken(token)
	if !valid {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.EndRequest()
		return
	}
	ctx.Params().Set("userID", userID)
	ctx.Params().Set("deviceID", deviceID)
	ctx.Next()
}

// ErrorHandler 统一错误处理
var ErrorHandler = func(ctx context.Context) {
	v := ctx.Values().Get(CtxKeyHandlerError)
	if v != nil {
		if err, ok := v.(error); ok {
			log.Printf("Handler error, Method: %s, Path: %s, Error: %v", ctx.Method(), ctx.Path(), err)
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.WriteString("Error Occurred: " + err.Error())
		}
	}
}
