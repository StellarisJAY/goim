package db

import "github.com/go-redis/redis/v8"

var DB *Databases

type Databases struct {
	MySQL *MysqlDB
	Redis *redis.Client
}

func init() {
	DB = new(Databases)
	sql, err := InitMySQL()
	if err != nil {
		panic(err)
	}
	DB.MySQL = sql
	DB.Redis = newRedisDB()
}
