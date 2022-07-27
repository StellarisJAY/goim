package dao

import (
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
)

func InsertMessage(msg *model.Message) error {
	tx := db.DB.MySQL.Create(msg)
	return tx.Error
}

// ListMessages 获取两个用户之间一段时间内的消息列表
func ListMessages(user1, user2 int64, startTime, endTime int64) ([]*model.Message, error) {
	messages := make([]*model.Message, 0)
	tx := db.DB.MySQL.
		Where("user1=? and user2=? and timestamp between ? and ?", user1, user2, startTime, endTime).
		Find(messages)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messages, nil
}
