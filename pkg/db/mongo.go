package db

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	MongoDBName               = "db_goim"
	CollectionOfflineMessage  = "offlineMessage"
	CollectionGroupInvitation = "groupInvitation"
	CollectionFriendRequest   = "friendRequest"
	CollectionNotification    = "notification"
)

var Day = int64(time.Hour) * 24

func InitMongoDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(&options.ClientOptions{
		Hosts: config.Config.MongoDB.Hosts,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	err = client.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("connect to mongoDB error %w", err))
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("ping mongoDB error %w", err))
	}

	// 创建离线消息表，离线消息只保留7天
	offlineMessageOptions := options.CreateCollection().SetExpireAfterSeconds(Day * 7)
	_ = client.Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionOfflineMessage, offlineMessageOptions)

	uniqueOption := options.Index().SetUnique(true)
	expireOption := options.CreateCollection().SetExpireAfterSeconds(Day * 3)
	// 创建通知表，通知仅保存三天
	_ = client.Database(MongoDBName).CreateCollection(context.TODO(), CollectionNotification, expireOption)
	// 创建通知的唯一索引，避免重复通知（通知类型仅有：好友申请、进群邀请相关的通知，这些通知不能重复出现）
	_, err = client.Database(MongoDBName).Collection(CollectionNotification).
		Indexes().
		CreateOne(context.TODO(),
			mongo.IndexModel{
				Keys:    bson.D{{"receiver", 1}, {"triggerUser", 1}, {"type", 1}},
				Options: uniqueOption,
			})
	if err != nil {
		panic(fmt.Errorf("create notification unique index error %w", err))
	}

	return client, nil
}
