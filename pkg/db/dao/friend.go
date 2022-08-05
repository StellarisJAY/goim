package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
)

const (
	KeyFriendRelation = "user_friends_"
)

func GetFriendInfo(userID, friendID int64) (*model.Friend, error) {
	// 缓存获取好友信息
	marshal, err := db.DB.Redis.HGet(context.TODO(), fmt.Sprintf("%s%d", KeyFriendRelation, userID), fmt.Sprintf("%x", friendID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，从MySQL读取
			friend := &model.Friend{}
			tx := db.DB.MySQL.
				Table("friends").
				Where("owner_id=? AND friend_id=?", userID, friendID).
				Take(friend)
			if tx.Error != nil {
				return nil, tx.Error
			}
			// 写回缓存
			_ = CacheFriendship(friend)
			return friend, nil
		}
	}
	friend := &model.Friend{}
	err = json.Unmarshal(marshal, friend)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func ListFriends(userID int64) ([]*model.Friend, error) {
	friends := make([]*model.Friend, 0)
	tx := db.DB.MySQL.
		Table("friends").
		Where("owner_id=?", userID).
		Find(&friends)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return friends, nil
}

func InsertFriendship(friendships ...*model.Friend) error {
	return db.DB.MySQL.CreateInBatches(friendships[:], 2).Error
}

func CacheFriendship(friendship *model.Friend) error {
	marshal, err := json.Marshal(friendship)
	if err != nil {
		return err
	}
	hSet := db.DB.Redis.HSet(context.TODO(), fmt.Sprintf("%s%d", KeyFriendRelation, friendship.OwnerID), fmt.Sprintf("%x", friendship.FriendID), marshal)
	return hSet.Err()
}
