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

func (s *Server) HandleNSQ(message *nsq.Message) error {
	userID, payload, err := splitNsqMessage(message)
	if err != nil {
		return fmt.Errorf("consume nsq message error: %w", err)
	}
	value, ok := s.wsServer.UserConns.Load(userID)
	if !ok {
		log.Debug("message to offline user: %d", userID)
		return nil
	}
	if config.Config.Gateway.UseJsonMsg {
		payload, err = protoMessageToJsonMessage(payload)
		if err != nil {
			return fmt.Errorf("can't transfer protobuf message to json: %w", err)
		}
	}
	channel := value.(*websocket.Channel)
	if err := channel.Push(payload); err != nil {
		return fmt.Errorf("push message to channel %s error %w", channel.ID(), err)
	}
	return nil
}

func (s *Server) HandleKafka(message *sarama.ConsumerMessage) error {
	userID, err := strconv.ParseInt(string(message.Key), 10, 64)
	if err != nil {
		return fmt.Errorf("kafka message corrupted, parse userID error %w", err)
	}
	value, ok := s.wsServer.UserConns.Load(userID)
	if !ok {
		log.Debug("message to offline user: %d", userID)
		return nil
	}
	channel := value.(*websocket.Channel)
	payload := message.Value
	if config.Config.Gateway.UseJsonMsg {
		if p, err := protoMessageToJsonMessage(payload); err != nil {
			return fmt.Errorf("transfer proto message to json error %w", err)
		} else {
			payload = p
		}
	}
	if err := channel.Push(payload); err != nil {
		return fmt.Errorf("push message to channel %s for user %d error %w", channel.ID(), userID, err)
	}
	return nil
}

func protoMessageToJsonMessage(payload []byte) ([]byte, error) {
	message := &pb.BaseMsg{}
	if err := proto.Unmarshal(payload, message); err != nil {
		return nil, err
	}
	return json.Marshal(message)
}

func splitNsqMessage(message *nsq.Message) (int64, []byte, error) {
	buffer := bytes.NewBuffer(message.Body)
	key, err := buffer.ReadBytes(';')
	if err != nil {
		return 0, nil, fmt.Errorf("split nsq message error %w", err)
	}
	key = key[:len(key)-1]
	userID, err := strconv.ParseInt(string(key), 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("parse userID from nsq message error %w", err)
	}
	return userID, buffer.Bytes(), nil
}
