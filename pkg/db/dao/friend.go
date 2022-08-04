package dao

import (
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
)

func GetFriendInfo(userID, friendID int64) (*model.Friend, error) {
	friend := &model.Friend{}
	tx := db.DB.MySQL.
		Table("friends").
		Where("owner_id=? AND friend_id=?", userID, friendID).
		Take(friend)
	if tx.Error != nil {
		return nil, tx.Error
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
