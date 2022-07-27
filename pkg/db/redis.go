package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/config"
)

func newRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:  "",
		Addr:     config.Config.Redis.Address,
		Username: config.Config.Redis.User,
		Password: config.Config.Redis.Password,
	})
}
