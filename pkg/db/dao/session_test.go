package dao

import (
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"testing"
)

func InitBatchGetGroupSessions() {
	_, _, _ = SaveSession(1001, "device-a", "gate-1", "channel-1")
	_, _, _ = SaveSession(1001, "device-b", "gate-1", "channel-2")
	_, _, _ = SaveSession(1002, "device-2a", "gate-1", "channel-3")
	_, _, _ = SaveSession(1002, "device-2b", "gate-2", "channel-1")
	_ = AddGroupMember(&model.GroupMember{GroupID: 10001, UserID: 1001})
	_ = AddGroupMember(&model.GroupMember{GroupID: 10001, UserID: 1002})
}

func TestBatchGetGroupSessions(t *testing.T) {
	InitBatchGetGroupSessions()
	t.Run("redis-lua-test", func(t *testing.T) {
		sessions, err := BatchGetGroupSessions(10001, "device-a", 1001)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		for k, v := range sessions {
			fmt.Printf("%d : %#v\n", k, v)
		}
	})
}

func BenchmarkBatchGetGroupSessions(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = BatchGetGroupSessions(10001, "device-a", 1001)
		}
	})
}

func BenchmarkGetGroupSessions(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = GetGroupSessions(10001, "device-a", 1001)
		}
	})
}
