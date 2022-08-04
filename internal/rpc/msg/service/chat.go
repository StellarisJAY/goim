package service

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"time"
)

func (m *MessageServiceImpl) SendMessage(ctx context.Context, request *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	message := request.Msg
	// 为消息添加时间戳 和 ID
	message.Timestamp = time.Now().UnixMilli()
	message.Id = m.idGenerator.NextID()
	// 序列化
	marshal, err := proto.Marshal(message)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 发送给推送服务
	key := fmt.Sprintf("%x", message.Id)
	_, _, err = m.transferProducer.PushMessage(key, marshal)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.SendMsgResponse{Code: pb.Success, MessageId: message.Id, Timestamp: message.Timestamp}, nil
}
