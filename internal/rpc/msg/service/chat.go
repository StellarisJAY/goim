package service

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/kafka"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"google.golang.org/protobuf/proto"
	"time"
)

type ChatServiceImpl struct {
	transferProducer *kafka.Producer
	idGenerator      *snowflake.Snowflake
}

func NewChatServiceImpl() *ChatServiceImpl {
	transProducer, err := kafka.NewProducer(config.Config.Kafka.Addrs, pb.MessageTransferTopic)
	if err != nil {
		panic(err)
	}
	return &ChatServiceImpl{
		transferProducer: transProducer,
		idGenerator:      snowflake.NewSnowflake(config.Config.MachineID),
	}
}

func (c *ChatServiceImpl) SendMessage(ctx context.Context, request *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	message := request.Msg
	// 为消息添加时间戳 和 ID
	message.Timestamp = time.Now().UnixMilli()
	message.Id = c.idGenerator.NextID()
	// 序列化
	marshal, err := proto.Marshal(message)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 发送给推送服务
	key := fmt.Sprintf("%x", message.Id)
	_, _, err = c.transferProducer.PushMessage(key, marshal)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.SendMsgResponse{Code: pb.Success, MessageId: message.Id, Timestamp: message.Timestamp}, nil
}