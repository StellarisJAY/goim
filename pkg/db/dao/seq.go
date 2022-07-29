package dao

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/db"
)

// 用户收件箱序列号用 Redis 的 Hash保存，结构如下
// user_seq_#userId: {
//     group_1: 0
//     group_2: 122
//     group_3: 22
//     non_group: 100
// }
// 序列号表示当前群聊中的最新消息序号，用户通过与本地记录的序号比较来判断是否同步新消息
//
const (
	userSeqKey  = "user_seq_%d"
	nonGroupKey = "0"
)

// IncrUserSeq 增加用户收件箱的 序列号
func IncrUserSeq(userId int64) (int64, error) {
	command := db.DB.Redis.HIncrBy(context.Background(), fmt.Sprintf(userSeqKey, userId), nonGroupKey, 1)
	return command.Result()
}

// IncrUserGroupChatSeq 增加某个用户在某个群聊中的收件序列号
func IncrUserGroupChatSeq(userId, groupId int64) (int64, error) {
	command := db.DB.Redis.HIncrBy(context.Background(), fmt.Sprintf(userSeqKey, userId), fmt.Sprintf("%d", groupId), 1)
	return command.Result()
}
