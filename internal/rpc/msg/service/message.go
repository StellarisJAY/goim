package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/kafka"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"github.com/stellarisJAY/goim/pkg/wordfilter"
)

type MessageServiceImpl struct {
	transferProducer *kafka.Producer
	idGenerator      *snowflake.Snowflake
	wordFilter       wordfilter.Filter
}

func NewMessageServiceImpl() *MessageServiceImpl {
	transProducer, err := kafka.NewProducer(config.Config.Kafka.Addrs, pb.MessageTransferTopic)
	if err != nil {
		panic(err)
	}
	// 从配置文件读取敏感词
	filter := wordfilter.NewTrieTreeFilter()
	filter.Build(config.Config.SensitiveWords)
	return &MessageServiceImpl{
		transferProducer: transProducer,
		idGenerator:      snowflake.NewSnowflake(config.Config.MachineID),
		wordFilter:       filter,
	}
}

// SyncOfflineMessages 同步离线消息
// 1. 从MongoDB查询用户提供的序列号开始的消息
// 2. 将消息按照序列号排序
// 3. 返回消息列表
func (m *MessageServiceImpl) SyncOfflineMessages(ctx context.Context, request *pb.SyncMsgRequest) (*pb.SyncMsgResponse, error) {
	var messages []*model.OfflineMessage
	var err error
	if *request.Flag == int32(pb.MessageFlag_Group) {
		messages, err = dao.ListOfflineGroupMessages(request.UserID, *request.From, request.LastSeq)
	} else {
		// 查询 to == userID AND seq > lastSeq
		messages, err = dao.ListOfflineMessages(request.UserID, request.LastSeq)
	}
	if err != nil {
		return &pb.SyncMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 转换message到baseMsg
	msgs := make([]*pb.BaseMsg, len(messages))
	for i, message := range messages {
		msgs[i] = &pb.BaseMsg{
			Id:        message.ID,
			From:      message.From,
			To:        message.To,
			Content:   string(message.Content),
			Flag:      pb.MessageFlag(message.Flag),
			Timestamp: message.Timestamp,
			Seq:       message.Seq,
		}
	}
	response := &pb.SyncMsgResponse{
		Code:     pb.Success,
		Message:  "",
		Messages: msgs,
	}
	// 更新客户端的 initSeq 和 lastSeq
	if msgs != nil && len(msgs) > 0 {
		response.InitSeq = msgs[0].Seq
		response.LastSeq = msgs[len(msgs)-1].Seq
	}
	return response, nil
}
