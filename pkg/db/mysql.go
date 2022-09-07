package db

import (
	"fmt"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type MysqlDB struct {
	sync.RWMutex
	*gorm.DB
}

func InitMySQL() (*MysqlDB, error) {
	DBName := config.Config.MySQL.DB
	connect := fmt.Sprintf("%s:%s@tcp(%s)/mysql?charset=utf8", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Address)
	db, err := gorm.Open(mysql.Open(connect))
	if err != nil {
		return nil, fmt.Errorf("can't connect to mysql database: %w", err)
	}
	// 创建数据库
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", DBName)
	err = db.Exec(createDB).Error
	if err != nil {
		return nil, fmt.Errorf("can't create database in target mysql server: %w", err)
	}

	sql := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Address, DBName)
	db, err = gorm.Open(mysql.Open(sql))
	if err != nil {
		return nil, err
	}
	// 创建表
	_ = db.AutoMigrate(&model.User{}, &model.DeviceLogin{})
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.Migrator().HasTable(&model.User{}) {
		_ = db.Migrator().CreateTable(&model.User{})
	}
	if !db.Migrator().HasTable(&model.DeviceLogin{}) {
		_ = db.Migrator().CreateTable(&model.DeviceLogin{})
	}
	if !db.Migrator().HasTable(&model.Message{}) {
		_ = db.Migrator().CreateTable(&model.Message{})
	}
	if !db.Migrator().HasTable(&model.Group{}) {
		_ = db.Migrator().CreateTable(&model.Group{})
	}
	if !db.Migrator().HasTable(&model.GroupMember{}) {
		_ = db.Migrator().CreateTable(&model.GroupMember{})
	}
	if !db.Migrator().HasTable(&model.Friend{}) {
		_ = db.Migrator().CreateTable(&model.Friend{})
	}
	return &MysqlDB{
		RWMutex: sync.RWMutex{},
		DB:      db,
	}, nil
}
