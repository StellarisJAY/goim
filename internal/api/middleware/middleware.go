package middleware

import (
	_context "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/authutil"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/stringutil"
	"time"
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
	userID, deviceID, valid, expired, expireAt := authutil.ValidateToken(token)
	if !valid {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.WriteString("invalid token")
		ctx.EndRequest()
		return
	}
	if expired {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.WriteString("token expired")
		ctx.EndRequest()
		return
	}
	// token有效时间小于等于10分钟，更新token
	if time.Now().Sub(time.UnixMilli(expireAt)).Minutes() <= 10 {
		authService, err := GetAuthService()
		if err != nil {
			log.Errorf("update token for user %s failed, error: %v", userID, err)
		} else {
			uid, _ := stringutil.HexStringToInt64(userID)
			response, err := authService.UpdateToken(_context.TODO(), &pb.UpdateTokenRequest{UserID: uid, DeviceID: deviceID})
			if err != nil || response.Code == pb.Error {
				log.Errorf("update token for user: %s failed, message: ", userID, response.Message)
			} else {
				recorder := ctx.Recorder()
				recorder.Header().Set("AuthUpdateToken", token)
			}
		}
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
			log.Errorf("Handler error, Method: %s, Path: %s, Error: %v", ctx.Method(), ctx.Path(), err)
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.WriteString("Internal Error: " + err.Error())
		} else if errMsg, ok := v.(string); ok {
			log.Errorf("Handler error, Method: %s, Path: %s, Error Msg: %s", ctx.Method(), ctx.Path(), errMsg)
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.WriteString("Internal Error: " + errMsg)
		}
	}
}

func GetAuthService() (pb.AuthClient, error) {
	// 从服务发现获取 RPC 客户端连接
	conn, err := naming.GetClientConn(pb.UserServiceName)
	if err != nil {
		return nil, err
	}
	// RPC调用用户注册服务
	return pb.NewAuthClient(conn), nil
}
