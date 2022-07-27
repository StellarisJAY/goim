package handler

import (
	_context "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"log"
)

var LoginHandler context.Handler = func(ctx context.Context) {
}

// RegisterHandler 用户注册API
var RegisterHandler context.Handler = func(ctx context.Context) {
	// 读取JSON请求
	regReq := new(pb.RegisterRequest)
	err := ctx.ReadJSON(regReq)
	if err != nil {
		// 请求解析失败，BadRequest
		_, _ = ctx.Problem(iris.NewProblem().Status(iris.StatusBadRequest))
		return
	}
	// 从服务发现获取 RPC 客户端连接
	conn, err := naming.GetClientConn("auth")
	if err != nil {
		_, _ = ctx.Problem(iris.NewProblem().Status(iris.StatusInternalServerError))
		return
	}
	// RPC调用用户注册服务
	client := pb.NewAuthClient(conn)
	response, err := client.Register(_context.Background(), regReq)
	if err != nil {
		log.Println("failed to register: ", err)
		_, _ = ctx.Problem(iris.NewProblem().Status(iris.StatusInternalServerError))
		return
	}
	_, _ = ctx.JSON(response)
}
