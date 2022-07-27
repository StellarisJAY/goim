package snowflake

import (
	"fmt"
	"sync"
	"time"
)

type Snowflake struct {
	sequence      int64 // 当前时间戳下的序列号
	lastTimestamp int64 // 上一个时间戳
	machineID     int64 // 机器ID
	mutex         sync.Mutex
}

const (
	epoch          int64 = 1645539742000 // 起始时间戳，从2022年2月2日22：22：22开始计时
	timestampBits        = 41
	machineIdBits        = 10
	sequenceBits         = 12
	timestampMask  int64 = 1<<timestampBits - 1
	timestampShift       = machineIdBits + sequenceBits
	sequenceMask   int64 = 0xfff
)

func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{
		sequence:      0,
		lastTimestamp: time.Now().UnixNano()/1000000 - epoch,
		machineID:     machineID << sequenceBits,
		mutex:         sync.Mutex{},
	}
}

// SetMachineID 初始化雪花算法时设置的机器ID
func (s *Snowflake) SetMachineID(mid int64) {
	s.machineID = mid << sequenceBits
}

// NextID 获取下一个SnowflakeID
func (s *Snowflake) NextID() int64 {
	timestamp := time.Now().UnixNano()/1000000 - epoch
	s.mutex.Lock()
	// 相同毫秒，seq递增
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		// seq溢出，等待下一毫秒
		if s.sequence == 0 {
			for timestamp == s.lastTimestamp {
				timestamp = time.Now().UnixNano()/1000000 - epoch
			}
			s.lastTimestamp = timestamp
		}
	} else {
		s.sequence = 0
		s.lastTimestamp = timestamp
	}
	seq := s.sequence
	s.mutex.Unlock()
	return seq | s.machineID | ((timestamp & timestampMask) << timestampShift)
}

func (s *Snowflake) NextHexString() string {
	return fmt.Sprintf("%x", s.NextID())
}
