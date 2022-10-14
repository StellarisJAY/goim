package dao

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/cache"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"gorm.io/gorm"
	"strconv"
)

const sessionPrefix = "goim_session_"

// 保存session的脚本
// 因为保存session需要检查该设备是否已经登录在某个网关，使用单个Redis命令无法做到。
// 所以需要脚本执行
var saveSessionScript = `
	local old = redis.call('hget', KEYS[1], ARGV[1])
	redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
	return old
`

// SaveSession 保存某个用户的某台设备的登录信息
// 如果查询到该设备已经存在登录信息，则需要返回原来所在的网关和channel
func SaveSession(userId int64, deviceId, gateway, channel string) (string, string, error) {
	key := fmt.Sprintf("%s%d", sessionPrefix, userId)
	encodedSession := encodeSession(gateway, channel)
	eval := db.DB.Redis.Eval(context.Background(), saveSessionScript, []string{key}, deviceId, encodedSession)
	if res, err := eval.Result(); err != nil {
		if err == redis.Nil {
			return "", "", nil
		}
		return "", "", err
	} else if res != nil {
		if oldSession, ok := res.([]byte); ok {
			oldGate, oldChan := decodeSession(oldSession)
			return oldGate, oldChan, nil
		}
	}
	return "", "", nil
}

// GetSessions 获取除了 fromDevice 以外 用户的所有登录设备 session 信息
func GetSessions(userId int64, fromDevice string, fromUser int64) ([]model.Session, error) {
	key := fmt.Sprintf("%s%d", sessionPrefix, userId)
	all := db.DB.Redis.HGetAll(context.Background(), key)
	if encodeds, err := all.Result(); err != nil {
		return nil, err
	} else {
		if fromUser == userId {
			// 排除来源设备
			delete(encodeds, fromDevice)
		}
		sessions := make([]model.Session, 0, len(encodeds))
		for _, encoded := range encodeds {
			gateway, channel := decodeSession([]byte(encoded))
			sessions = append(sessions, model.Session{Gateway: gateway, Channel: channel})
		}
		return sessions, nil
	}
}

func GetGroupSessions(groupId int64, fromDevice string, fromUser int64) (map[int64][]model.Session, error) {
	// 获取群成员IDs
	groupMembers, err := ListGroupMemberIDs(groupId)
	if err != nil {
		return nil, err
	}
	sessions := make(map[int64][]model.Session)
	for _, userID := range groupMembers {
		if err != nil {
			continue
		}
		session, err := GetSessions(userID, fromDevice, fromUser)
		if err != nil {
			continue
		}
		sessions[userID] = session
	}
	return sessions, nil
}

func KickSession(userID int64, deviceID string) error {
	key := fmt.Sprintf("%s%d", sessionPrefix, userID)
	del := db.DB.Redis.HDel(context.TODO(), key, deviceID)
	return del.Err()
}

func GroupMemberExists(groupID int64, userID int64) (bool, error) {
	return cache.IsMember(fmt.Sprintf(KeyGroupMembers, groupID), strconv.FormatInt(userID, 10), func(g string, m string) (bool, error) {
		member := &model.GroupMember{}
		result := db.DB.MySQL.
			Table("group_members").
			Where("group_id=? AND user_id=?", groupID, userID).
			First(member)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return false, result.Error
		} else if result.Error != nil {
			return false, nil
		} else {
			return true, nil
		}
	})
}

// session 编码格式：4字节gateLen + 4字节chanLen + gateway + channel
func encodeSession(gateway, channel string) []byte {
	data := make([]byte, len(gateway)+len(channel)+8)
	binary.BigEndian.PutUint32(data[0:4], uint32(len(gateway)))
	binary.BigEndian.PutUint32(data[4:8], uint32(len(channel)))
	copy(data[8:], gateway)
	copy(data[8+len(gateway):], channel)
	return data
}

func decodeSession(data []byte) (gateway, channel string) {
	lenGate := binary.BigEndian.Uint32(data[0:4])
	lenChan := binary.BigEndian.Uint32(data[4:8])
	gateway = string(data[8 : 8+lenGate])
	channel = string(data[8+lenGate : 8+lenGate+lenChan])
	return
}
