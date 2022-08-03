package model

const (
	MessageFlagFrom  byte = 0
	MessageFlagTo    byte = 1
	MessageFlagGroup byte = 2

	MemberStatusNormal  byte = 0
	MemberStatusInvited byte = 1
	MemberStatusBanned  byte = 2

	MemberRoleOwner byte = iota
	MemberRoleAdmin
	MemberRoleNormal
)

// User 用户表
type User struct {
	ID           int64  `gorm:"column:id;primaryKey;type:int8"`
	Account      string `gorm:"column:account;unique;type:varchar(255)"`
	Password     string `gorm:"column:password;type:varchar(255)"`
	NickName     string `gorm:"column:nick_name;type:varchar(64)"`
	Salt         string `gorm:"column:salt"`
	RegisterTime int64  `gorm:"column:register_time;type:int8"`
}

type UserInfo struct {
	ID           int64  `json:"id"`
	Account      string `json:"account"`
	NickName     string `json:"nickName"`
	RegisterTime int64  `json:"registerTime"`
}

// DeviceLogin 设备登录记录表
type DeviceLogin struct {
	UserID    int64  `gorm:"column:user_id;type:int8"`
	DeviceID  string `gorm:"column:device_id;type:varchar(255)"`
	Timestamp int64  `gorm:"column:timestamp;type:int8"`
	Ip        string `gorm:"column:ip;type:varchar(16)"`
}

// Message 持久化消息表
type Message struct {
	ID        int64  `gorm:"column:id;type:int8;primaryKey"`
	User1     int64  `gorm:"column:user1;type:int8;index:idx_user_timestamp"`
	User2     int64  `gorm:"column:user2;type:int8;index:idx_user_timestamp"`
	Content   []byte `gorm:"column:content;type:varchar(255)"`
	Timestamp int64  `gorm:"column:timestamp;type:int8;index:idx_user_timestamp"`
	Flag      byte   `gorm:"column:flag;type:int1"`
}

// OfflineMessage 离线消息表
type OfflineMessage struct {
	From      int64  `json:"from" bson:"from"`
	To        int64  `json:"to" bson:"to"`
	Content   []byte `json:"content" bson:"content"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Seq       int64  `json:"seq" bson:"seq"`   // 序列号为接收用户的自增序列号，用户通过本地的序列号和消息序列号判断是否同步消息
	Flag      byte   `json:"flag" bson:"flag"` // Flag 标记消息目标的类型
}

// Session 信息
type Session struct {
	Gateway string // 网关地址
	Channel string // 网关内部的channel
}

// Group 群组表
type Group struct {
	ID           int64  `gorm:"column:id;type:int8;primaryKey" json:"id"`
	Name         string `gorm:"column:name;type:varchar(64);" json:"name"`
	CreateTime   int64  `gorm:"column:create;type:int8" json:"createTime"`
	Description  string `gorm:"column:description;type:varchar(255)" json:"description"`
	OwnerID      int64  `gorm:"column:owner_id;type:int8;" json:"ownerID"`
	OwnerAccount string `gorm:"column:owner_account;type:varchar(255)" json:"ownerAccount"`
}

// GroupMember 群成员表
type GroupMember struct {
	GroupID  int64 `gorm:"column:group_id;type:int8;primaryKey"`
	UserID   int64 `gorm:"column:user_id;type:int8;primaryKey"`
	JoinTime int64 `gorm:"column:join_time;type:int8"`
	Status   byte  `gorm:"column:status;type:tinyint"` // 群成员状态：正常、已邀请未加入、禁言
	Role     byte  `gorm:"column:role;type:tinyint"`   // 群成员角色：群主、管理员、普通成员
}

// GroupMemberFull 群成员详细信息model
type GroupMemberFull struct {
	*GroupMember
	Account  string `gorm:"column:account"`
	NickName string `gorm:"column:nick_name"`
}

// GroupInvitation 进群邀请记录，仅在MongoDB保存三天
type GroupInvitation struct {
	ID             int64  `bson:"id" json:"id""`
	UserID         int64  `bson:"userID" json:"userID"`
	GroupID        int64  `bson:"groupID" json:"groupID"`
	Timestamp      int64  `bson:"timestamp" json:"timestamp"`
	Inviter        int64  `bson:"inviter" json:"inviter"`
	InviterAccount string `json:"inviterAccount"`
	GroupName      string `json:"groupName" gorm:"column:name"` // 不在MongoDB保存
}

// Friend 好友关系表
type Friend struct {
	User1      int64 `gorm:"column:user1;type:int8;primaryKey"`
	User2      int64 `gorm:"column:user2;type:int8;primaryKey"`
	AcceptTime int64 `gorm:"column:accept_time;type:int8;"`
}

// AddFriendRequest 好友请求
type AddFriendRequest struct {
	Requester int64  `bson:"requester"`
	Target    int64  `bson:"target"`
	Timestamp int64  `bson:"timestamp"`
	Message   string `bson:"message"` // 验证信息
}
