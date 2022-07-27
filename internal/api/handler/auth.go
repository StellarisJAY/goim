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
)

var validate = validator.New()

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {}, http.AuthRequest{})
}

// AuthHandler 授权用户设备，返回访问Token
var AuthHandler context.Handler = func(ctx context.Context) {
	req := new(http.AuthRequest)
	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	authReq := &pb.AuthRequest{}
	if err := copier.CopyStructFields(authReq, req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	// 获取授权RPC服务，RPC调用获取Token
	service, err := GetAuthService()
	if err != nil {
		_, _ = ctx.Problem(iris.StatusInternalServerError)
		return
	}
	response, err := service.AuthorizeDevice(_context.Background(), authReq)
	if err != nil {
		_, _ = ctx.Problem(iris.StatusInternalServerError)
	}
	_, _ = ctx.JSON(&http.AuthResponse{
		BaseResponse: http.BaseResponse{
			Code: response.Code, Message: response.Message,
		},
		Token: response.Token,
	})
}

// RegisterHandler 用户注册API
var RegisterHandler context.Handler = func(ctx context.Context) {
	// 读取JSON请求
	regReq := new(http.RegisterRequest)
	if err := ctx.ReadJSON(regReq); err != nil {
		_, _ = ctx.Problem(iris.NewProblem().Status(iris.StatusBadRequest))
		return
	}
	if err := validate.Struct(regReq); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	request := &pb.RegisterRequest{}
	if err := copier.CopyStructFields(request, regReq); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	service, err := GetAuthService()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	if response, err := service.Register(_context.Background(), request); err != nil {
		_, _ = ctx.Problem(iris.NewProblem().Status(iris.StatusInternalServerError))
		return
	} else {
		_, _ = ctx.JSON(&http.BaseResponse{
			Code:    response.Code,
			Message: response.Message,
		})
	}
}

func GetAuthService() (pb.AuthClient, error) {
	// 从服务发现获取 RPC 客户端连接
	conn, err := naming.GetClientConn("auth")
	if err != nil {
		return nil, err
	}
	// RPC调用用户注册服务
	return pb.NewAuthClient(conn), nil
}
