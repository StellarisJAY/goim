package model

import "time"

// User 用户表
type User struct {
	ID        int64     `gorm:"column:id;primary key"`
	Account   string    `gorm:"column:account;unique"`
	Password  string    `gorm:"column:password"`
	NickName  string    `gorm:"column:nick_name"`
	CreatedAt time.Time `gorm:"column:create_time"`
}

// DeviceLogin 设备登录记录表
type DeviceLogin struct {
	UserID    int64  `gorm:"column:user_id"`
	DeviceID  string `gorm:"column:device_id"`
	Timestamp int64  `gorm:"column:timestamp"`
	Ip        string `gorm:"column:ip"`
}
