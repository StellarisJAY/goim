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

type SyncOfflineGroupMessageRequest struct {
	Groups     []int64 `json:"groups" validate:"required"`
	Timestamps []int64 `json:"timestamps" validate:"required"`
}

type SyncOfflineGroupMessageResponse struct {
	BaseResponse
	GroupMessages []*pb.SingleGroupMessages
}

type SyncGroupMessageResponse struct {
	BaseResponse
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

// GroupInfoResponse 查询群聊信息返回
type GroupInfoResponse struct {
	BaseResponse
	Group *pb.GroupInfo `json:"group"`
}

// FindUserResponse 查询用户信息返回
type FindUserResponse struct {
	BaseResponse
	UserInfo *pb.UserInfo `json:"userInfo"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	NickName string `json:"nickName"`
}

// ListGroupMemberResponse 列出群成员返回
type ListGroupMemberResponse struct {
	BaseResponse
	Members []*pb.GroupMember `json:"members"`
}

// AddFriendRequest 添加好友请求
type AddFriendRequest struct {
	TargetID        int64  `json:"targetID" validate:"required"`
	ValidateMessage string `json:"validateMessage" validate:"required"`
}

// AcceptFriendRequest 接受好友申请请求
type AcceptFriendRequest struct {
	TargetID int64 `json:"targetID" validate:"required"`
}

type ListFriendsResponse struct {
	BaseResponse
	Friends []*pb.FriendInfo `json:"friends"`
}

type GetFriendInfoResponse struct {
	BaseResponse
	Info *pb.FriendInfo `json:"info"`
}

type SyncGroupLatestMessagesRequest struct {
	GroupID int64 `json:"groupID" validate:"required"`
	Limit   int64 `json:"limit" validate:"required"`
}

type SyncGroupLatestMessagesResponse struct {
	BaseResponse
	FirstSeq int64         `json:"firstSeq"`
	LastSeq  int64         `json:"lastSeq"`
	Msgs     []*pb.BaseMsg `json:"msgs"`
}

type ListJoinedGroupsResponse struct {
	BaseResponse
	Groups []*pb.GroupInfo `json:"groups"`
}

type ListNotificationRequest struct {
	BaseResponse
	Notifications []*Notification
}

type Notification struct {
	Id          int64  `json:"id"`
	Receiver    int64  `json:"receiver"`    // 通知接收者
	TriggerUser int64  `json:"triggerUser"` // 通知触发者
	Type        byte   `json:"type"`        // 通知类型
	Message     string `json:"message"`     // 内容
	Read        bool   `json:"read"`
	Timestamp   int64  `json:"timestamp"`
}
