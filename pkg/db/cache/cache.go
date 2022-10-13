package cache

import (
	"context"
	"encoding/json"
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

func Get(key string, expire time.Duration, missFunc func(key string) (interface{}, error)) (interface{}, error) {
	res := db.DB.Redis.Get(context.TODO(), key)
	if res.Err() != nil && res.Err() == redis.Nil {
		if value, err := missFunc(key); err != nil {
			return nil, err
		} else if value != nil {
			marshal, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			db.DB.Redis.Set(context.TODO(), key, marshal, expire)
			return value, nil
		}
		return nil, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	return []byte(res.Val()), nil
}

func Delete(key string) error {
	del := db.DB.Redis.Del(context.TODO(), key)
	return del.Err()
}

func IsMember(key, member string, missFunc func(string, string) (bool, error)) (bool, error) {
	if res := db.DB.Redis.SIsMember(context.TODO(), key, member); res.Err() != nil && res.Err() == redis.Nil {
		isMember, err := missFunc(key, member)
		if err != nil {
			return false, err
		}
		if isMember {
			db.DB.Redis.SAdd(context.TODO(), key, member)
		}
		return isMember, nil
	} else if res.Err() != nil {
		return false, res.Err()
	} else {
		return res.Val(), nil
	}
}
