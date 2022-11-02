package api

import (
	"github.com/kataras/iris/v12"
	"github.com/stellarisJAY/goim/internal/api/handler"
	"github.com/stellarisJAY/goim/internal/api/middleware"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
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
		messageParty.Get("/offline/group", handler.SyncOfflineGroupMessages)
		messageParty.Get("/offline/group/latest", handler.SyncLatestGroupMessages)
	}

	notificationParty := application.Party("/notification")
	{
		notificationParty.Use(middleware.TokenVerifier)
		notificationParty.Done(middleware.ErrorHandler)
		notificationParty.Get("/list", handler.ListNotifications)
		notificationParty.Post("/mark/read/{id:int64}", handler.MarkNotificationReadHandler)
	}
	// 用户信息API
	userParty := application.Party("/user")
	{
		userParty.Use(middleware.TokenVerifier)
		userParty.Done(middleware.ErrorHandler)
		userParty.Get("/{id:int64}", handler.FindUserHandler)
		userParty.Put("", handler.UpdateUserHandler)
	}
	// 群聊相关API
	groupParty := application.Party("/group")
	{
		groupParty.Use(middleware.TokenVerifier)
		groupParty.Done(middleware.ErrorHandler)
		groupParty.Post("", handler.CreateGroupHandler)
		groupParty.Get("/{id:int64}", handler.GroupInfoHandler)
		groupParty.Get("/members/{id:int64}", handler.GroupMemberHandler)
		groupParty.Post("/{gid:int64}/invite/{uid:int64}", handler.InviteUserHandler)
		groupParty.Get("/invitations", handler.ListInvitationsHandler)
		groupParty.Post("/invitation/{invID:int64}", handler.AcceptInvitationHandler)
		groupParty.Get("/list", handler.ListJoinedGroupsHandler)
	}
	// 好友相关API
	friendParty := application.Party("/friend")
	{
		friendParty.Use(middleware.TokenVerifier)
		friendParty.Done(middleware.ErrorHandler)
		// 添加好友申请
		friendParty.Post("/application", handler.InsertFriendApplicationHandler)
		// 列出好友申请
		//friendParty.Get("/application", handler.ListFriendApplications)
		// 接受好友申请
		friendParty.Put("/application", handler.AcceptFriendHandler)
		// 获取好友信息
		friendParty.Get("/{friendID:int64}", handler.FriendInfoHandler)
		// 获取好友列表
		friendParty.Get("/list", handler.FriendListHandler)
	}
}

func Start() {
	err := application.Run(iris.Addr(":" + config.Config.ApiServer.Port))
	if err != nil {
		log.Errorf("iris http server error: %v", err)
	}
}
