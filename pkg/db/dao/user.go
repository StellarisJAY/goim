package dao

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/cache"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

const (
	KeyUserInfo = "user_info_%d"
)

func FindUserByAccount(account string) (*model.User, bool, error) {
	user := new(model.User)
	tx := db.DB.MySQL.Where("account=?", account).First(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, tx.Error
	}
	return user, true, nil
}

func FindUserInfo(userID int64) (*model.UserInfo, error) {
	key := fmt.Sprintf(KeyUserInfo, userID)
	result, err := cache.Get(key, 0, func(key string) (*model.User, error) {
		user := &model.User{}
		if tx := db.DB.MySQL.Model(user).Where("id=?", userID).First(user); tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		} else if tx.Error != nil {
			return nil, tx.Error
		} else {
			return user, nil
		}
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &model.UserInfo{ID: result.ID, Account: result.Account, NickName: result.NickName, RegisterTime: result.RegisterTime}, nil
}

func UpdateUserNickname(user *model.UserInfo) error {
	// 更新数据库
	tx := db.DB.MySQL.Where("id=?", user.ID).Find(&model.User{}).UpdateColumn("nick_name", user.NickName)
	if tx.Error != nil {
		return tx.Error
	}
	// 删除缓存内容
	_ = cache.Delete(fmt.Sprintf(KeyUserInfo, user.ID))
	return nil
}

func InsertUser(user *model.User) error {
	tx := db.DB.MySQL.Create(user)
	return tx.Error
}

func InsertUserLoginLog(login *model.DeviceLogin) error {
	tx := db.DB.MySQL.Create(login)
	return tx.Error
}

// AddFriendRequest 添加好友请求
func AddFriendRequest(request *model.AddFriendRequest) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	// 写入到MongoDB, 返回错误可能是重复添加
	_, err := database.Collection(db.CollectionFriendRequest).InsertOne(context.TODO(), request)
	return err
}

// GetAndDeleteFriendRequest 获取并删除好友请求，删除成功表示接收添加好友请求
func GetAndDeleteFriendRequest(requester, target int64) (*model.AddFriendRequest, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	request := new(model.AddFriendRequest)
	err := database.Collection(db.CollectionFriendRequest).FindOneAndDelete(context.TODO(), bson.D{
		{"target", target},
		{"requester", requester},
	}).Decode(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

// ListAddFriendRequests 列出添加当前用户为好友的所有请求
func ListAddFriendRequests(target int64) ([]*model.AddFriendRequest, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	cursor, err := database.Collection(db.CollectionFriendRequest).Find(context.TODO(), bson.D{
		{"target", target},
	})
	if err != nil {
		return nil, err
	}
	requests := make([]*model.AddFriendRequest, 0)
	err = cursor.All(context.TODO(), &requests)
	if err != nil {
		return nil, err
	}
	return requests, nil
}
