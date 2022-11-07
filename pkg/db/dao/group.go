package dao

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/cache"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/stringutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

const (
	KeyGroupInfo    = "group_info_%d"
	KeyJoinedGroup  = "user_joined_group_%d"
	KeyGroupMembers = "group_members_%d"
)

func InsertGroup(group *model.Group) error {
	tx := db.DB.MySQL.Create(group)
	return tx.Error
}

func FindGroupInfo(groupID int64) (*model.Group, error) {
	key := fmt.Sprintf(KeyGroupInfo, groupID)
	res, err := cache.Get(key, 0, func(key string) (*model.Group, error) {
		group := &model.Group{}
		tx := db.DB.MySQL.Table("groups").Where("id=?", groupID).First(group)
		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		} else if tx.Error != nil {
			return nil, tx.Error
		} else {
			return group, nil
		}
	})
	if err != nil {
		return nil, err
	}
	return res, nil
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

func ListGroupMemberIDs(groupID int64) ([]int64, error) {
	memberIDs, err := cache.ListMembers(fmt.Sprintf(KeyGroupMembers, groupID), 0, func(key string) ([]string, error) {
		var memberIDs []int64
		result := db.DB.MySQL.
			Select("user_id").
			Table("group_members").
			Where("group_id=?", groupID).
			Find(&memberIDs)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		} else if result.Error != nil {
			return nil, nil
		} else {
			return stringutil.Int64ListToString(memberIDs), nil
		}
	})
	if err != nil {
		return nil, err
	}
	return stringutil.StringListToInt64(memberIDs), nil
}

func ListStringGroupMemberIDs(groupID int64) ([]string, error) {
	memberIDs, err := cache.ListMembers(fmt.Sprintf(KeyGroupMembers, groupID), 0, func(key string) ([]string, error) {
		var memberIDs []int64
		result := db.DB.MySQL.
			Select("user_id").
			Table("group_members").
			Where("group_id=?", groupID).
			Find(&memberIDs)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		} else if result.Error != nil {
			return nil, nil
		} else {
			return stringutil.Int64ListToString(memberIDs), nil
		}
	})
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}

// FindGroupMember 查询群成员信息
func FindGroupMember(groupID, userID int64) (*model.GroupMember, error) {
	member := &model.GroupMember{}
	result := db.DB.MySQL.
		Table("group_members").
		Where("group_id=? AND user_id=?", groupID, userID).
		First(member)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	} else if result.Error != nil {
		return nil, nil
	} else {
		return member, nil
	}
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

// ListUserJoinedGroupIds 列出用户加入的所有群聊的ID
func ListUserJoinedGroupIds(userID int64) ([]int64, error) {
	key := fmt.Sprintf(KeyJoinedGroup, userID)
	groupIDs, err := cache.ListMembers(key, 0, func(key string) ([]string, error) {
		var groups []int64
		result := db.DB.MySQL.
			Select("group_id").
			Table("group_members").
			Where("user_id=?", userID).
			Find(&groups)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		} else if result.Error != nil || len(groups) == 0 {
			return nil, nil
		} else {
			return stringutil.Int64ListToString(groups), nil
		}
	})
	if err != nil {
		return nil, err
	}
	return stringutil.StringListToInt64(groupIDs), nil
}

// ListGroupInfos 列出给定ID的所有群聊基本信息
func ListGroupInfos(groupIDs []int64) ([]*model.Group, error) {
	groups := make([]*model.Group, 0)
	result := db.DB.MySQL.
		Select("*").
		Table("groups").
		Where("id in ?", groupIDs).
		Find(&groups)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return groups, nil
}

func RemoveGroupMember(groupID int64, memberID int64) error {
	tx := db.DB.MySQL.Unscoped().Delete(&model.GroupMember{GroupID: groupID, UserID: memberID})
	if tx.Error != nil {
		return tx.Error
	}
	_ = cache.Delete(fmt.Sprintf(KeyGroupMembers, groupID))
	_ = cache.Delete(fmt.Sprintf(KeyJoinedGroup, memberID))
	return nil
}
