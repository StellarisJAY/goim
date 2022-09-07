package dao

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertMessage(msg *model.Message) error {
	tx := db.DB.MySQL.Create(msg)
	return tx.Error
}

// InsertMessages 批量插入消息
func InsertMessages(messages []*model.Message) error {
	tx := db.DB.MySQL.CreateInBatches(messages, len(messages))
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

func ListOfflineMessages(userID int64, lastSeq int64) ([]*model.OfflineMessage, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	// 按照seq排序
	opts := options.Find().SetSort(bson.D{{"seq", 1}})
	query := bson.D{
		{"to", userID},
		{"seq", bson.D{{"$gt", lastSeq}}},
	}
	result, err := database.Collection(db.CollectionOfflineMessage).Find(context.TODO(), query, opts)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.OfflineMessage, 0)
	if result.All(context.TODO(), &messages) != nil {
		return nil, err
	}
	return messages, nil
}

func ListOfflineGroupMessages(userID int64, groupID int64, lastTimestamp int64) ([]*model.OfflineMessage, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	// 按照seq排序
	opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
	query := bson.D{
		{"Flag", pb.MessageFlag_Group},
		{"to", groupID},
		{"timestamp", bson.D{{"$gt", lastTimestamp}}},
	}
	result, err := database.Collection(db.CollectionOfflineMessage).Find(context.TODO(), query, opts)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.OfflineMessage, 0)
	if err = result.All(context.TODO(), &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// ListLatestOfflineGroupMessages 查询群聊最新的离线消息
func ListLatestOfflineGroupMessages(groupID int64, lastTimestamp int64, limit int32) ([]*model.OfflineMessage, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	// 按照timestamp 倒序
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(int64(limit))
	var query bson.D
	// lastTimeStamp 为 -1，从最新一条消息开始查询
	if lastTimestamp == -1 {
		query = bson.D{
			{"Flag", pb.MessageFlag_Group},
			{"to", groupID},
		}
	} else {
		query = bson.D{
			{"Flag", pb.MessageFlag_Group},
			{"to", groupID},
			{"timestamp", bson.D{{"$lt", lastTimestamp}}},
		}
	}
	cur, err := database.Collection(db.CollectionOfflineMessage).Find(context.TODO(), query, opts)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.OfflineMessage, limit)
	if err = cur.All(context.TODO(), &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
