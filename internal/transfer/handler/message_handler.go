package handler

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/pool"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"runtime"
)

var pushWorker = pool.NewWorkerPool(runtime.NumCPU() * 2)

// MessageTransferHandler 消息中转处理器
// 1. 从消息队列消费到一条消息后，
var MessageTransferHandler = func(message *sarama.ConsumerMessage) error {
	value := message.Value
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	var err error
	switch byte(msg.Flag) {
	case pb.MessageFlagFrom:
		err = handleSingleMessage(msg)
	case pb.MessageFlagGroup:
		handleGroupChat(msg)
	default:
	}
	if err != nil {
		log.Println(err)
	}
	return nil
}

// handleSingleMessage 单聊消息处理
// 1. 在Redis自增用户的收件序列号，然后用自增的结果作为当前消息的序号
// 2. 获取到序列号后，将消息写入MongoDB作为离线消息保存。
// 3. Mongo写入完成后，尝试向目标用户所在的网关发送RPC推送消息
func handleSingleMessage(message *pb.BaseMsg) error {
	seq, err := dao.IncrUserSeq(message.To)
	if err != nil {
		return fmt.Errorf("increment user sequence number error %w", err)
	}
	message.Seq = seq
	offlineMessage := &model.OfflineMessage{
		From:      message.From,
		To:        message.To,
		Content:   []byte(message.Content),
		Timestamp: message.Timestamp,
		Seq:       seq,
		Flag:      pb.MessageFlagFrom,
	}
	err = dao.InsertOfflineMessage(offlineMessage)
	if err != nil {
		return fmt.Errorf("insert offline message error %w", err)
	}
	return pushMessage(message)
}

// handleGroupChat 群聊消息处理
func handleGroupChat(message *pb.BaseMsg) {

}

func pushMessage(message *pb.BaseMsg) error {
	// 查询目标用户所在的session
	sessions, err := dao.GetSessions(message.To, message.DeviceId)
	log.Println("user session: ", sessions)
	if err != nil {
		return fmt.Errorf("get session info error: %w", err)
	}
	for _, session := range sessions {
		pushWorker.Submit(func() {
			// 与gateway连接
			conn, psErr := naming.DialConnection(session.Gateway)
			if psErr != nil {
				log.Println("connect gateway error: ", psErr)
				return
			}
			// RPC push message
			client := pb.NewRelayClient(conn)
			response, psErr := client.PushMessage(context.Background(), &pb.PushMsgRequest{Message: message, Channel: session.Channel})
			if psErr != nil {
				log.Println("push message RPC error: ", psErr)
				return
			}
			if response.Base.Code != pb.Success {
				log.Println("push message result: ", response.Base.Message)
			}
		})
	}
	return nil
}
