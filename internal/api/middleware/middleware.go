package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/authutil"
	"log"
)

var TokenVerifier = func(ctx context.Context) {
	var token string
	if token = ctx.GetHeader("Authorization"); token == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.EndRequest()
		return
	}
	log.Println("token: ", token)
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
