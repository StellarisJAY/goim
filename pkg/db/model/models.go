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
	User1     int64  `gorm:"column:user1;type:int8;primaryKey"`
	User2     int64  `gorm:"column:user2;type:int8;primaryKey"`
	Content   []byte `gorm:"column:content;type:varchar(255)"`
	Timestamp int64  `gorm:"column:timestamp;type:int8;primaryKey"`
	Flag      byte   `gorm:"column:flag;type:int1"`
}

type OfflineMessage struct {
}
