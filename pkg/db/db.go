package db

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *Databases

type Databases struct {
	MySQL   *MysqlDB
	Redis   *redis.Client
	MongoDB *mongo.Client
}

func init() {
	DB = new(Databases)
	sql, err := InitMySQL()
	if err != nil {
		panic(err)
	}
	DB.MySQL = sql
	DB.Redis = newRedisDB()
	DB.MongoDB, err = InitMongoDB()
	if err != nil {
		panic(err)
	}
}
