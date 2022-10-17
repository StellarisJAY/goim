package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/mq"
	"github.com/stellarisJAY/goim/pkg/mq/kafka"
	_nsq "github.com/stellarisJAY/goim/pkg/mq/nsq"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/pool"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"runtime"
	"strconv"
	"strings"
)

var (
	pushWorker            = pool.NewWorkerPool(runtime.NumCPU() * 2)
	onlineMessageProducer mq.MessageProducer
)

func init() {
	configMQ := strings.ToLower(config.Config.MessageQueue)
	switch configMQ {
	case "nsq":
		onlineMessageProducer = _nsq.NewProducer()
	case "kafka":
		producer, err := kafka.NewProducer(config.Config.Kafka.Addrs, pb.MessagePushTopic)
		if err != nil {
			panic(err)
		}
		onlineMessageProducer = producer
	default:
		panic(fmt.Errorf("unsupported message queue: %s", configMQ))
	}
}

// MessageTransferHandler 消息中转处理器
// 1. 从消息队列消费到一条消息后，
var MessageTransferHandler = func(message *sarama.ConsumerMessage) error {
	value := message.Value
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	var err error
	switch msg.Flag {
	case pb.MessageFlag_From:
		err = handleSingleMessage(msg)
	case pb.MessageFlag_Group:
		err = handleGroupChat(msg)
	default:
	}
	if err != nil {
		log.Warn("handle message failed, msgID: %d, error: %v", msg.Id, err)
	}
	return nil
}

var NsqMessageHandler = func(message *nsq.Message) error {
	value := message.Body
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	var err error
	switch msg.Flag {
	case pb.MessageFlag_From:
		err = handleSingleMessage(msg)
	case pb.MessageFlag_Group:
		err = handleGroupChat(msg)
	default:
	}
	if err != nil {
		log.Warn("handle message failed, msgID: %d, error: %v", msg.Id, err)
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
		ID:        message.Id,
		From:      message.From,
		To:        message.To,
		Content:   []byte(message.Content),
		Timestamp: message.Timestamp,
		Seq:       seq,
		Flag:      byte(message.Flag),
	}
	err = dao.InsertOfflineMessage(offlineMessage)
	if err != nil {
		return fmt.Errorf("insert offline message error %w", err)
	}
	return pushMessage(message)
}

// handleGroupChat 群聊消息处理
func handleGroupChat(message *pb.BaseMsg) error {
	offlineMessage := &model.OfflineMessage{
		ID:        message.Id,
		From:      message.From,
		To:        message.To,
		Content:   []byte(message.Content),
		Timestamp: message.Timestamp,
		Seq:       0,
		Flag:      byte(pb.MessageFlag_Group),
	}
	err := dao.InsertOfflineMessage(offlineMessage)
	if err != nil {
		return fmt.Errorf("insert offline message error %w", err)
	}
	return pushGroupMessage(message)
}

// pushMessage 推送消息，分为同步和异步两种推送方式
// 同步推送调用Gateway的RPC服务发送给用户所处的指定网关，该过程中需要查询Redis缓存寻址
// 异步推送直接把消息发送给消息队列，由Gateway异步消费后发送给客户端
func pushMessage(message *pb.BaseMsg) error {
	syncPush := config.Config.SyncPushOnline
	if syncPush {
		return pushOnlineSync(message)
	} else {
		return pushOnlineMQ(message)
	}
}

// pushGroupMessage 推送群聊消息，尝试向群聊中的每个用户推送消息
func pushGroupMessage(message *pb.BaseMsg) error {
	sessions, err := dao.GetGroupSessions(message.To, message.DeviceId, message.From)
	if err != nil {
		return fmt.Errorf("get group session error: %w", err)
	}
	// 通过每个用户的sessions，整理出每个Gateway需要推送消息的channels
	gates := make(map[string][]string)
	for _, v := range sessions {
		for _, session := range v {
			channels, ok := gates[session.Gateway]
			if !ok {
				channels = make([]string, 0)
			}
			channels = append(channels, session.Channel)
			gates[session.Gateway] = channels
		}
	}
	// 向每个网关推送消息
	for gateway, channels := range gates {
		pushWorker.Submit(func() {
			conn, err := naming.DialConnection(gateway)
			if err != nil {
				log.Warn("Connect to gateway failed, gateway: ", gateway)
				return
			}
			client := pb.NewRelayClient(conn)
			// 向网关服务发送广播消息请求
			_, err = client.Broadcast(context.TODO(), &pb.BroadcastRequest{
				Message:  message,
				Channels: channels,
			})
			if err != nil {
				log.Warn("RPC Broadcast message failed, gateway: %s, error: %v", gateway, err)
			}
		})
	}
	return nil
}

// pushOnlineSync 同步发送
// 1. 从Redis寻址，获取用户所在的网关器服务地址
// 2. 通过RPC向网关发送消息
func pushOnlineSync(message *pb.BaseMsg) error {
	// 查询目标用户所在的session
	sessions, err := dao.GetSessions(message.To, message.DeviceId, message.From)
	if err != nil {
		return fmt.Errorf("get session info error: %w", err)
	}
	for _, session := range sessions {
		pushWorker.Submit(func() {
			// 与gateway连接
			conn, psErr := naming.DialConnection(session.Gateway)
			if psErr != nil {
				log.Warn("connect gateway error: ", psErr)
				return
			}
			// RPC push message
			client := pb.NewRelayClient(conn)
			response, psErr := client.PushMessage(context.Background(), &pb.PushMsgRequest{Message: message, Channel: session.Channel})
			if psErr != nil {
				log.Warn("push message RPC error: ", psErr)
				return
			}
			if response.Base.Code != pb.Success {
				log.Info("push message result: ", response.Base.Message)
			}
		})
	}
	return nil
}

func pushOnlineMQ(message *pb.BaseMsg) error {
	if marshal, err := proto.Marshal(message); err != nil {
		return fmt.Errorf("push message marshal error: %w", err)
	} else {
		key := strconv.FormatInt(message.To, 10)
		switch onlineMessageProducer.(type) {
		case *_nsq.Producer:
			body := bytes.Buffer{}
			body.Write([]byte(key))
			body.Write([]byte(";"))
			body.Write(marshal)
			return onlineMessageProducer.PushMessage(pb.MessagePushTopic, key, body.Bytes())
		default:
			return onlineMessageProducer.PushMessage(pb.MessagePushTopic, key, marshal)
		}
	}
}
