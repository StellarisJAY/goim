package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/copier"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"gorm.io/gorm"
)

const (
	KeyUserInfo = "user_info_"
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
	userInfo := new(model.UserInfo)
	// 查询Redis
	res := db.DB.Redis.Get(context.TODO(), fmt.Sprintf("%s%d", KeyUserInfo, userID))
	marshal, err := res.Bytes()
	if err != nil {
		if err == redis.Nil {
			user := &model.User{}
			// Redis中没有，查询MySQL
			tx := db.DB.MySQL.Where("id=?", userID).Find(user)
			if tx.Error != nil {
				return nil, err
			}
			_ = copier.CopyStructFields(userInfo, user)
			// 写入Redis
			marshal, _ = json.Marshal(userInfo)
			_ = db.DB.Redis.Set(context.TODO(), fmt.Sprintf("%s%d", KeyUserInfo, userID), marshal, 0)
			return userInfo, nil
		}
		return nil, err
	}
	err = json.Unmarshal(marshal, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func UpdateUserInfo(user *model.UserInfo) error {
	// 更新数据库
	tx := db.DB.MySQL.Where("id=?", user.ID).Find(&model.User{}).UpdateColumn("nick_name", user.NickName)
	if tx.Error != nil {
		return tx.Error
	}
	// 删除缓存内容
	_ = db.DB.Redis.Del(context.TODO(), fmt.Sprintf("%s%d", KeyUserInfo, user.ID))
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
