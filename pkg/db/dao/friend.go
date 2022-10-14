package dao

import (
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/cache"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/stringutil"
	"gorm.io/gorm"
	"strconv"
)

const (
	KeyFriendIDList = "user_friend_ids_%d"
)

func CheckFriendship(userID int64, friendID int64) (bool, error) {
	key := fmt.Sprintf(KeyFriendIDList, userID)
	return cache.IsMember(key, strconv.FormatInt(friendID, 10), func(key string, member string) (bool, error) {
		friend := &model.Friend{}
		tx := db.DB.MySQL.Table("friends").
			Where("owner_id=? AND friend_id=?", userID, friendID).
			First(friend)
		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			return false, nil
		} else if tx.Error != nil {
			return false, tx.Error
		} else {
			return true, nil
		}
	})
}

func ListFriendIDs(userID int64) ([]int64, error) {
	members, err := cache.ListMembers(fmt.Sprintf(KeyFriendIDList, userID), 0, func(key string) ([]string, error) {
		var friends []int64
		tx := db.DB.MySQL.Table("friends").
			Where("owner_id=?", userID).
			Select("friend_id").
			Find(&friends)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return stringutil.Int64ListToString(friends), nil
	})
	if err != nil {
		return nil, err
	}
	return stringutil.StringListToInt64(members), nil
}

func InsertFriendship(friendships ...*model.Friend) error {
	return db.DB.MySQL.CreateInBatches(friendships[:], 2).Error
}
