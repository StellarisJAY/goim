package dao

import (
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
)

func InsertGroup(group *model.Group) error {
	tx := db.DB.MySQL.Create(group)
	return tx.Error
}

func InsertGroupMember(groupMember *model.GroupMember) error {
	tx := db.DB.MySQL.Create(groupMember)
	return tx.Error
}

func ListGroupMembers(groupID int64) ([]*model.GroupMember, error) {
	members := make([]*model.GroupMember, 0)
	tx := db.DB.MySQL.Where("group_id=?", groupID).Find(members)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return members, nil
}
