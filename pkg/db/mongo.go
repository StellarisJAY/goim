package db

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InitMongoDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(&options.ClientOptions{
		Hosts: config.Config.MongoDB.Hosts,
	})
	if err != nil {
		return nil, err
	}

	database := client.Database("goim")
	err = database.CreateCollection(context.TODO(), "offline_messages")
	if err != nil {
		log.Println("create offline messages collection error: ", err)
	}
	return client, nil
}
