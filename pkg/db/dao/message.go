package dao

import (
	"context"
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

// InsertOfflineMessage 保存一条离线消息
func InsertOfflineMessage(msg *model.OfflineMessage) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionOfflineMessage)
	_, err := collection.InsertOne(context.Background(), msg, nil)
	return err
}

// InsertOfflineMessages 批量保存离线消息
func InsertOfflineMessages(messages []*model.OfflineMessage) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionOfflineMessage)
	temp := make([]interface{}, len(messages))
	for i, msg := range messages {
		temp[i] = msg
	}
	_, err := collection.InsertMany(context.Background(), temp, nil)
	return err
}
