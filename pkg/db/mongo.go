package db

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	MongoDBName               = "db_goim"
	CollectionOfflineMessage  = "offlineMessage"
	CollectionGroupInvitation = "groupInvitation"
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
	offlineMessageExpire := Day * int64(config.Config.Message.OfflineExpireTime)
	_ = client.
		Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionOfflineMessage,
			&options.CreateCollectionOptions{ExpireAfterSeconds: &offlineMessageExpire})
	groupInvitationExpire := Day * 3
	_ = client.Database(MongoDBName).
		CreateCollection(context.TODO(), CollectionGroupInvitation, &options.CreateCollectionOptions{ExpireAfterSeconds: &groupInvitationExpire})
	return client, nil
}
