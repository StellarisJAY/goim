package model

const (
	MessageFlagFrom  byte = 0
	MessageFlagTo    byte = 1
	MessageFlagGroup byte = 2
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
