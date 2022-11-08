package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/panjf2000/ants/v2"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/mq"
	"github.com/stellarisJAY/goim/pkg/mq/kafka"
	_nsq "github.com/stellarisJAY/goim/pkg/mq/nsq"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"runtime"
	"strconv"
	"strings"
)

var (
	pushWorkerPool        *ants.Pool
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
	if wp, err := ants.NewPool(runtime.NumCPU() * 2); err != nil {
		panic(fmt.Errorf("error occurred when creating push worker: %w", err))
	} else {
		pushWorkerPool = wp
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
		log.Warn("handle message failed", zap.Int64("messageID", msg.Id), zap.Error(err))
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
		log.Warn("handle message failed", zap.Int64("messageID", msg.Id), zap.Error(err))
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
	// 同步模式使用RPC广播给gateway服务
	if config.Config.Transfer.SyncPushOnline {
		return pushGroupSync(message)
	} else {
		// 异步模式使用MQ转发消息
		return pushGroupMQ(message)
	}
}

// pushMessage 推送消息，分为同步和异步两种推送方式
// 同步推送调用Gateway的RPC服务发送给用户所处的指定网关，该过程中需要查询Redis缓存寻址
// 异步推送直接把消息发送给消息队列，由Gateway异步消费后发送给客户端
func pushMessage(message *pb.BaseMsg) error {
	syncPush := config.Config.Transfer.SyncPushOnline
	if syncPush {
		return pushOnlineSync(message)
	} else {
		return pushOnlineMQ(message)
	}
}

// pushGroupSync 推送群聊消息，尝试向群聊中的每个用户推送消息
func pushGroupSync(message *pb.BaseMsg) error {
	sessions, err := dao.BatchGetGroupSessions(message.To, message.DeviceId, message.From)
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
		_ = pushWorkerPool.Submit(func() {
			groupPushTask(gateway, channels, message)
		})
	}
	return nil
}

func groupPushTask(gateway string, channels []string, message *pb.BaseMsg) {
	conn, err := naming.DialConnection(gateway)
	if err != nil {
		log.Warn("Connect to gateway failed ", zap.String("gateway", gateway), zap.Error(err))
		return
	}
	client := pb.NewRelayClient(conn)
	// 向网关服务发送广播消息请求
	_, err = client.Broadcast(context.TODO(), &pb.BroadcastRequest{
		Message:  message,
		Channels: channels,
	})
	if err != nil {
		log.Warn("RPC Broadcast message failed", zap.String("gateway", gateway), zap.Error(err))
	}
}

func syncPushTask(session model.Session, message *pb.BaseMsg) {
	// 与gateway连接
	conn, psErr := naming.DialConnection(session.Gateway)
	if psErr != nil {
		log.Warn("connect to gateway failed", zap.Error(psErr))
		return
	}
	// RPC push message
	client := pb.NewRelayClient(conn)
	response, psErr := client.PushMessage(context.Background(), &pb.PushMsgRequest{Message: message, Channel: session.Channel})
	if psErr != nil {
		log.Warn("push message RPC failed", zap.Error(psErr))
		return
	}
	if response.Base.Code != pb.Success {
		log.Info("push message finished", zap.String("result", response.Base.Message))
	}
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
		_ = pushWorkerPool.Submit(func() {
			syncPushTask(session, message)
		})
	}
	return nil
}

func pushOnlineMQ(message *pb.BaseMsg) error {
	key := strconv.FormatInt(message.To, 10)
	return pushToMQ(message, key)
}

func pushGroupMQ(message *pb.BaseMsg) error {
	groupID := message.To
	groupMembers, err := dao.ListStringGroupMemberIDs(groupID)
	if err != nil {
		return err
	}
	mqMsg := messageToMqGroupMessage(message)
	mqMsg.GroupMembers = groupMembers
	err = pushToMQ(mqMsg, "group")
	if err != nil {
		log.Warn("push group message to mq failed",
			zap.Int64("groupID", groupID),
			zap.Int64("messageID", message.Id),
			zap.Error(err))
	}
	return nil
}

// pushToMQ 消息推送到消息队列
func pushToMQ(message proto.Message, key string) error {
	if marshal, err := proto.Marshal(message); err != nil {
		return fmt.Errorf("push message marshal error: %w", err)
	} else {
		switch onlineMessageProducer.(type) {
		case *_nsq.Producer:
			// 由于NSQ消息队列没有消息key，所以需要在消息体中手动插入key
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

func messageToMqGroupMessage(baseMsg *pb.BaseMsg) *pb.MqGroupMessage {
	return &pb.MqGroupMessage{
		From:      baseMsg.From,
		To:        baseMsg.To,
		Content:   baseMsg.Content,
		Flag:      baseMsg.Flag,
		Timestamp: baseMsg.Timestamp,
		Id:        baseMsg.From,
		Seq:       baseMsg.Seq,
		DeviceId:  baseMsg.DeviceId,
	}
}
