package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/db"
	"time"
)

func List(key string, expire time.Duration, missFunc func(key string) []string) ([]string, error) {
	res := db.DB.Redis.LRange(context.TODO(), key, 0, -1)
	if res.Err() != nil && res.Err() == redis.Nil {
		values := missFunc(key)
		if values != nil && len(values) > 0 {
			db.DB.Redis.RPush(context.TODO(), key, values)
		}
		return values, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	result, _ := res.Result()
	return result, nil
}

func ListMembers(key string, expire time.Duration, missFunc func(key string) []string) ([]string, error) {
	res := db.DB.Redis.SMembers(context.TODO(), key)
	if res.Err() != nil && res.Err() == redis.Nil {
		values := missFunc(key)
		if values != nil && len(values) > 0 {
			db.DB.Redis.SAdd(context.TODO(), key, values)
		}
		return values, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	result, _ := res.Result()
	return result, nil
}

func Get(key string, expire time.Duration, missFunc func(key string) []byte) ([]byte, error) {
	res := db.DB.Redis.Get(context.TODO(), key)
	if res.Err() != nil && res.Err() == redis.Nil {
		value := missFunc(key)
		if value != nil && len(value) > 0 {
			db.DB.Redis.Set(context.TODO(), key, value, 0)
		}
		return value, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	result, _ := res.Result()
	return []byte(result), nil
}

func Delete(key string) error {
	del := db.DB.Redis.Del(context.TODO(), key)
	return del.Err()
}
