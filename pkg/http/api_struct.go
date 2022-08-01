package http

import "github.com/stellarisJAY/goim/pkg/proto/pb"

type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	Account  string `json:"account" validate:"required"`
	NickName string `json:"nickName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	BaseResponse
}

type AuthRequest struct {
	Account  string `json:"account" validate:"required"`
	DeviceID string `json:"deviceID" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	BaseResponse
	Token string `json:"token"`
}

type SendMessageRequest struct {
	To      int64  `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
	Flag    int32  `json:"flag" validate:"required"`
}

type SendMessageResponse struct {
	BaseResponse
	MessageID int64 `json:"messageID"`
	Timestamp int64 `json:"timestamp"`
}

// SyncOfflineMessageResponse 同步离线消息返回
type SyncOfflineMessageResponse struct {
	BaseResponse
	// 同步后的第一个序列号
	InitSeq int64 `json:"initSeq"`
	// 同步后的最大序列号
	LastSeq int64 `json:"lastSeq"`
	// 消息列表
	Messages []*pb.BaseMsg `json:"messages"`
}

// CreateGroupRequest 创建群聊请求
type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateGroupResponse struct {
	BaseResponse
	GroupID int64 `json:"groupID"`
}

type GroupInfoResponse struct {
	BaseResponse
	Group *pb.GroupInfo `json:"group"`
}

type FindUserResponse struct {
	BaseResponse
	UserInfo *pb.UserInfo `json:"userInfo"`
}

type UpdateUserRequest struct {
	NickName string `json:"nickName"`
}
