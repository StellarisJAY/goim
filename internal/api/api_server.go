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

	// 授权服务API
	authParty := application.Party("/auth")
	{
		authParty.Done(middleware.ErrorHandler)
		authParty.Post("/login", handler.AuthHandler)
		authParty.Post("/register", handler.RegisterHandler)
	}
	// 聊天服务API
	chatParty := application.Party("/chat")
	{
		chatParty.Use(middleware.TokenVerifier)
		chatParty.Done(middleware.ErrorHandler)
		chatParty.Post("/send", handler.SendMessageHandler)
	}
	// 消息查询服务API
	messageParty := application.Party("/message")
	{
		messageParty.Use(middleware.TokenVerifier)
		messageParty.Done(middleware.ErrorHandler)
		messageParty.Get("/offline/{seq:int64}", handler.SyncOfflineMessageHandler)
	}
	// 用户信息API
	userParty := application.Party("/user")
	{
		userParty.Use(middleware.TokenVerifier)
		userParty.Done(middleware.ErrorHandler)
		userParty.Get("/{id:int64}", handler.FindUserHandler)
		userParty.Put("", handler.UpdateUserHandler)
	}
}

func Start() {
	err := application.Run(iris.Addr(":" + config.Config.ApiServer.Port))
	if err != nil {
		log.Println(err)
	}
}
