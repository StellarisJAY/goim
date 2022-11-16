package handler

import (
	_context "context"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stellarisJAY/goim/pkg/http"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

var validate = validator.New()

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, http.AuthRequest{})
}

// LoginHandler 授权用户设备，返回访问Token
var LoginHandler context.Handler = func(ctx *context.Context) {
	req := new(http.AuthRequest)
	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	authReq := &pb.AuthRequest{
		Account:  req.Account,
		DeviceID: req.DeviceID,
		Password: req.Password,
	}
	// 获取授权RPC服务，RPC调用获取Token
	service, err := GetAuthService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	response, err := service.AuthorizeDevice(_context.Background(), authReq)
	if err != nil {
		handleError(ctx, err)
		return
	}
	_ = ctx.JSON(&http.AuthResponse{
		BaseResponse: http.BaseResponse{
			Code: response.Code, Message: response.Message,
		},
		Token: response.Token,
	})
}

// RegisterHandler 用户注册API
var RegisterHandler context.Handler = func(ctx *context.Context) {
	// 读取JSON请求
	regReq := new(http.RegisterRequest)
	if err := ctx.ReadJSON(regReq); err != nil {
		_ = ctx.Problem(iris.NewProblem().Status(iris.StatusBadRequest))
		return
	}
	if err := validate.Struct(regReq); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	request := &pb.RegisterRequest{
		Account: regReq.Account, NickName: regReq.NickName, Password: regReq.Password,
	}
	service, err := GetAuthService()
	if err != nil {
		handleError(ctx, err)
		return
	}
	if response, err := service.Register(_context.Background(), request); err != nil {
		handleError(ctx, err)
		return
	} else {
		_ = ctx.JSON(&http.BaseResponse{
			Code:    response.Code,
			Message: response.Message,
		})
	}
}
