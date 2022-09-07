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
	"go.mongodb.org/mongo-driver/mongo"
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
	// 查询 to == userID AND seq > lastSeq
	messages, err = dao.ListOfflineMessages(request.UserID, request.LastSeq)
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

// SyncOfflineGroupMessages 同步离线的群聊消息
// 1. 遍历需要同步的群聊列表。
// 2. 查询每个群聊的离线消息。
// 3. 将离线消息从 DO 转换成 DTO
func (m *MessageServiceImpl) SyncOfflineGroupMessages(ctx context.Context, request *pb.SyncGroupMsgRequest) (*pb.SyncGroupMsgResponse, error) {
	groups, timestamps := request.Groups, request.Timestamps
	groupMessages := make([]*pb.SingleGroupMessages, len(groups))
	for i := 0; i < len(groups) && i < len(timestamps); i++ {
		groupID, lastTimestamp := groups[i], timestamps[i]
		messages, err := dao.ListOfflineGroupMessages(request.UserID, groupID, lastTimestamp)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			} else {
				return &pb.SyncGroupMsgResponse{Code: pb.Error, Message: err.Error()}, nil
			}
		}
		groupMessages[i] = &pb.SingleGroupMessages{
			GroupID:      groupID,
			StartTimeout: lastTimestamp,
			LastTimeout:  messages[len(messages)-1].Timestamp,
			Msgs:         OfflineMessagesToBaseMessages(messages),
		}
	}
	return &pb.SyncGroupMsgResponse{
		Code:          pb.Success,
		GroupMessages: groupMessages,
	}, nil
}

// SyncGroupLatestMessages 同步群聊中最近的 n 条消息
// 1. 从离线消息表倒序查询，如果request中的timestamp为-1，则从最新一条消息开始，否则从timestamp位置开始
// 2. 仅查询request中给出的limit条消息
func (m *MessageServiceImpl) SyncGroupLatestMessages(ctx context.Context, request *pb.SyncGroupLatestMessagesRequest) (*pb.SyncGroupLatestMessagesResponse, error) {
	messages, err := dao.ListLatestOfflineGroupMessages(request.GroupID, request.LastTimestamp, request.Limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.SyncGroupLatestMessagesResponse{Code: pb.Success}, nil
		}
		return &pb.SyncGroupLatestMessagesResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.SyncGroupLatestMessagesResponse{
		Code: pb.Success,
		Msgs: OfflineMessagesToBaseMessages(messages),
	}, nil
}

func OfflineMessagesToBaseMessages(offlineMessages []*model.OfflineMessage) []*pb.BaseMsg {
	msgs := make([]*pb.BaseMsg, len(offlineMessages))
	for i, m := range offlineMessages {
		msgs[i] = &pb.BaseMsg{
			From:      m.From,
			To:        m.From,
			Content:   string(m.Content),
			Flag:      pb.MessageFlag(m.Flag),
			Timestamp: m.Timestamp,
			Id:        m.ID,
		}
	}
	return msgs
}
