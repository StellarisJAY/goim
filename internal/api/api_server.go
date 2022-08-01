package api

import (
	"github.com/kataras/iris/v12"
	"github.com/stellarisJAY/goim/internal/api/handler"
	"github.com/stellarisJAY/goim/internal/api/middleware"
	"github.com/stellarisJAY/goim/pkg/config"
	"log"
)

var application *iris.Application

func Init() {
	application = iris.New()
	// 错误处理
	application.OnErrorCode(iris.StatusInternalServerError, handler.InternalError)
	application.OnErrorCode(iris.StatusNotFound, handler.NotFound)
	application.OnErrorCode(iris.StatusBadRequest, handler.BadRequest)

	// 授权服务API
	authParty := application.Party("/auth")
	{
		authParty.Post("/login", handler.AuthHandler)
		authParty.Post("/register", handler.RegisterHandler)
	}
	// 聊天服务API
	chatParty := application.Party("/chat")
	{
		chatParty.Use(middleware.TokenVerifier)
		chatParty.Post("/send", handler.SendMessageHandler)
	}
	// 消息查询服务API
	messageParty := application.Party("/message")
	{
		messageParty.Use(middleware.TokenVerifier)
		messageParty.Get("/offline/{seq:int64}", handler.SyncOfflineMessageHandler)
	}
}

func Start() {
	err := application.Run(iris.Addr(":" + config.Config.ApiServer.Port))
	if err != nil {
		log.Println(err)
	}
}
