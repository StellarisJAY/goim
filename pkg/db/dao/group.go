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

// InsertGroupInvitation 邀请用户进群记录
func InsertGroupInvitation(invitation *model.GroupInvitation) error {
	// 保存在MongoDB中，只保留3天
	database := db.DB.MongoDB.Database(db.MongoDBName)
	_, err := database.Collection(db.CollectionGroupInvitation).InsertOne(context.TODO(), invitation)
	return err
}

// ListInvitations 查询某个用户的被邀请记录
func ListInvitations(userID int64) ([]*model.GroupInvitation, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	findOpts := options.Find().
		SetSort(bson.D{{"timestamp", 1}}).
		SetShowRecordID(true)
	cursor, err := database.
		Collection(db.CollectionGroupInvitation).
		Find(context.TODO(), bson.D{{"userID", userID}}, findOpts)
	if err != nil {
		return nil, err
	}
	result := make([]*model.GroupInvitation, 0)
	err = cursor.All(context.TODO(), &result)
	return result, err
}

func ListGroupMembers(groupID int64) ([]*model.GroupMemberFull, error) {
	members := make([]*model.GroupMemberFull, 0)
	tx := db.DB.MySQL.
		Select([]string{"`group_id`", "`user_id`", "`join_time`", "`status`", "`account`", "`nick_name`"}).
		Table("group_members").
		Joins("inner join users on users.id=group_members.user_id").
		Where("group_id=?", groupID).
		Find(&members)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return members, nil
}

// FindGroupMember 查询群成员信息
func FindGroupMember(groupID, userID int64) *model.GroupMember {
	member := &model.GroupMember{}
	tx := db.DB.MySQL.
		Where("group_id=? AND user_id=?", groupID, userID).
		Find(member)
	if tx.Error != nil || member.GroupID != groupID {
		return nil
	}
	return member
}

// FindGroupMemberFull 查询群成员详细信息
func FindGroupMemberFull(groupID, userID int64) *model.GroupMemberFull {
	member := &model.GroupMemberFull{}
	tx := db.DB.MySQL.
		Select([]string{"`group_id`", "`user_id`", "`join_time`", "`status`", "`account`", "`nick_name`"}).
		Table("group_members").
		Joins("inner join users on users.id=group_members.user_id").
		Where("group_id=? AND user_id=?", groupID, userID).
		Find(&member)
	if tx.Error != nil || member.GroupID != groupID {
		return nil
	}
	return member
}

// FindGroupNames 查询目标groups的名称
func FindGroupNames(groups []int64) ([]string, error) {
	names := make([]string, 0)
	tx := db.DB.MySQL.
		Select("`name`").
		Table("groups").
		Where("id in ?", groups).
		Find(&names)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return names, nil
}

// GetAndDeleteInvitation  获取并删除进群邀请记录
func GetAndDeleteInvitation(invID int64) (*model.GroupInvitation, error) {
	invitation := &model.GroupInvitation{}
	database := db.DB.MongoDB.Database(db.MongoDBName)
	err := database.Collection(db.CollectionGroupInvitation).
		FindOneAndDelete(context.TODO(), bson.D{{"id", invID}}).
		Decode(invitation)
	if err != nil {
		return nil, err
	}
	return invitation, nil
}
