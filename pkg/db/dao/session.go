package dao

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
)

const sessionPrefix = "goim_session_"

func SaveSession(userId int64, deviceId, gateway, channel string) error {
	key := fmt.Sprintf("%s%d", sessionPrefix, userId)
	db.DB.Redis.HSet(context.Background(), key, deviceId, encodeSession(gateway, channel))
	return nil
}

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
