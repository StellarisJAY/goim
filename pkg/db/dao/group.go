package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	KeyGroupInfo = "group_info_"
)

func InsertGroup(group *model.Group) error {
	tx := db.DB.MySQL.Create(group)
	return tx.Error
}

func FindGroupInfo(groupID int64) (*model.Group, error) {
	group := &model.Group{}
	// 从Redis读取群信息
	marshal, err := db.DB.Redis.Get(context.TODO(), fmt.Sprintf("%s%d", KeyGroupInfo, groupID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 从MySQL获取
			tx := db.DB.MySQL.Where("id=?", groupID).Find(group)
			if tx.Error != nil {
				return nil, tx.Error
			}
			marshal, _ = json.Marshal(group)
			// 写入Redis
			_ = db.DB.Redis.Set(context.TODO(), fmt.Sprintf("%s%d", KeyGroupInfo, groupID), marshal, 0)
			return group, nil
		}
	}
	err = json.Unmarshal(marshal, group)
	return group, err
}

func AddGroupMember(groupMember *model.GroupMember) error {
	tx := db.DB.MySQL.Create(groupMember)
	return tx.Error
}

// InviteGroupMember 邀请用户进群记录
func InviteGroupMember(userID, groupID int64) error {
	// 保存在MongoDB中，只保留3天
	database := db.DB.MongoDB.Database(db.MongoDBName)
	_, err := database.Collection(db.CollectionGroupInvitation).InsertOne(context.TODO(), &model.GroupInvitation{
		UserID:    userID,
		GroupID:   groupID,
		Timestamp: time.Now().UnixMilli(),
	})
	return err
}

// ListInvitations 查询某个用户的被邀请记录
func ListInvitations(userID int64) ([]*model.GroupInvitation, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	findOpts := options.Find().SetSort(bson.D{{"timestamp", 1}})
	cursor, err := database.Collection(db.CollectionGroupInvitation).Find(context.TODO(), bson.D{{"userID", userID}}, findOpts)
	if err != nil {
		return nil, err
	}
	result := make([]*model.GroupInvitation, 0)
	err = cursor.All(context.TODO(), result)
	return result, err
}

func ListGroupMembers(groupID int64) ([]*model.GroupMember, error) {
	members := make([]*model.GroupMember, 0)
	tx := db.DB.MySQL.Where("group_id=?", groupID).Find(members)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return members, nil
}
