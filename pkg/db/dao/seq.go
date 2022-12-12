package dao

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
)

// 用户收件箱序列号用 Redis 的 Hash保存，结构如下
// user_seq_{userID}
// 序列号表示当前群聊中的最新消息序号，用户通过与本地记录的序号比较来判断是否同步新消息
// 每个群聊拥有单独的收件序号，格式：group_seq_{groupID}
const (
	userSeqKey  = "user_seq_%d"
	groupSeqKey = "group_seq_%d"
)

// IncrUserSeq 增加用户收件箱的 序列号
func IncrUserSeq(userId int64) (int64, error) {
	command := db.DB.Redis.IncrBy(context.Background(), fmt.Sprintf(userSeqKey, userId), 1)
	return command.Result()
}

// IncrGroupChatSeq 增加群聊的收件序列号
func IncrGroupChatSeq(groupId int64) (int64, error) {
	command := db.DB.Redis.IncrBy(context.Background(), fmt.Sprintf(groupSeqKey, groupId), 1)
	return command.Result()
}
