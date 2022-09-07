package db

import (
	"context"
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
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	// 创建离线消息缓存表，离线消息只保留7天
	offlineMessageOptions := options.CreateCollection().SetExpireAfterSeconds(Day * 7)
	_ = client.Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionOfflineMessage, offlineMessageOptions)

	uniqueOption := options.Index().SetUnique(true)
	// 创建群邀请表，设置邀请3天过期
	groupOptions := options.CreateCollection().SetExpireAfterSeconds(Day * 3)
	_ = client.Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionGroupInvitation, groupOptions)
	// 创建群邀请表的唯一索引
	_, err = client.Database(MongoDBName).
		Collection(CollectionGroupInvitation).
		Indexes().
		CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{"userID", 1}, {"groupID", 1}}, Options: uniqueOption})
	if err != nil {
		panic("failed to create collection: group_invitations, " + err.Error())
	}
	// 创建好友请求表，设置请求3天过期
	friendOptions := options.CreateCollection().SetExpireAfterSeconds(Day * 3)
	_ = client.Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionFriendRequest, friendOptions)
	// 创建好友请求表的唯一索引
	_, err = client.Database(MongoDBName).
		Collection(CollectionFriendRequest).
		Indexes().
		CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{"requester", 1}, {"target", 1}}, Options: uniqueOption})
	if err != nil {
		panic("failed to create collection: friend_requests, " + err.Error())
	}
	return client, nil
}
