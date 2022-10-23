package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/protobuf/proto"
	"strconv"
)

// HandleNSQ 处理NSQ消息队列消息
func (s *Server) HandleNSQ(message *nsq.Message) error {
	key, payload, err := splitNsqMessage(message)
	if err != nil {
		return fmt.Errorf("consume nsq message error: %w", err)
	}
	switch key {
	case "group":
		return relayGroupMessage(s, payload)
	default:
		// JSON化处理
		if config.Config.Gateway.UseJsonMsg {
			p, err := protoMessageToJsonMessage(payload)
			if err != nil {
				return fmt.Errorf("can't transfer protobuf message to json: %w", err)
			} else {
				payload = p
			}
		}
		return relayMqMessage(s, key, payload)
	}
}

// HandleKafka 处理kafka消息队列消息
func (s *Server) HandleKafka(message *sarama.ConsumerMessage) error {
	key := string(message.Key)
	switch key {
	case "group":
		return relayGroupMessage(s, message.Value)
	default:
		payload := message.Value
		if config.Config.Gateway.UseJsonMsg {
			p, err := protoMessageToJsonMessage(payload)
			if err != nil {
				return fmt.Errorf("can't transfer protobuf message to json: %w", err)
			} else {
				payload = p
			}
		}
		return relayMqMessage(s, key, payload)
	}
}

func protoMessageToJsonMessage(payload []byte) ([]byte, error) {
	message := &pb.BaseMsg{}
	if err := proto.Unmarshal(payload, message); err != nil {
		return nil, err
	}
	return json.Marshal(message)
}

func splitNsqMessage(message *nsq.Message) (string, []byte, error) {
	buffer := bytes.NewBuffer(message.Body)
	key, err := buffer.ReadBytes(';')
	if err != nil {
		return "", nil, fmt.Errorf("split nsq message error %w", err)
	}
	key = key[:len(key)-1]

	return string(key), buffer.Bytes(), nil
}

// relayGroupMessage 转发群聊消息
// *Server：Gateway 服务器实例
// payload []byte: 转发消息体，已经是json或proto格式
func relayGroupMessage(s *Server, payload []byte) error {
	mqMsg := &pb.MqGroupMessage{}
	err := proto.Unmarshal(payload, mqMsg)
	if err != nil {
		return fmt.Errorf("nsq message damaged %w", err)
	}
	baseMsg := mqMessageToBaseMsg(mqMsg)
	marshal, err := marshalBaseMsg(baseMsg)
	if err != nil {
		return fmt.Errorf("marshal base msg error %w", err)
	}
	for _, member := range mqMsg.GroupMembers {
		if err := relayMqMessage(s, member, marshal); err != nil {
			log.Debug("relay group member %s message error %w", member, err)
		}
	}
	return nil
}

// relayMqMessage 转发从消息队列收到的单聊消息
// *Server: gateway服务器实例
// key string： message key，一般是转发目标用户ID
// payload []byte: 转发消息体，已经转化为json或proto
func relayMqMessage(s *Server, key string, payload []byte) error {
	userID, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid userID format: %s, parse error: %w", key, err)
	}
	value, ok := s.wsServer.UserConns.Load(userID)
	if !ok {
		log.Debug("user not on gateway server: %d", userID)
		return nil
	}
	channel := value.(*websocket.Channel)
	if err := channel.Push(payload); err != nil {
		return fmt.Errorf("push message to channel %s error %w", channel.ID(), err)
	}
	return nil
}

// marshalBaseMsg 将基础的msg实例序列化成json或proto格式
func marshalBaseMsg(baseMsg *pb.BaseMsg) ([]byte, error) {
	if config.Config.Gateway.UseJsonMsg {
		if marshal, err := json.Marshal(baseMsg); err != nil {
			return nil, err
		} else {
			return marshal, nil
		}
	} else {
		return proto.Marshal(baseMsg)
	}
}

func mqMessageToBaseMsg(mqMsg *pb.MqGroupMessage) *pb.BaseMsg {
	return &pb.BaseMsg{
		From:      mqMsg.From,
		To:        mqMsg.To,
		Content:   mqMsg.Content,
		Flag:      mqMsg.Flag,
		Timestamp: mqMsg.Timestamp,
		Id:        mqMsg.Id,
		Seq:       mqMsg.Seq,
		DeviceId:  mqMsg.DeviceId,
	}
}
