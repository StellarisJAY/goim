package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/protobuf/proto"
	"strconv"
)

func (s *Server) HandleNSQ(message *nsq.Message) error {
	buffer := bytes.NewBuffer(message.Body)
	key, err := buffer.ReadBytes(';')
	if err != nil {
		return fmt.Errorf("nsq consumer read message key error %w", err)
	}
	key = key[:len(key)-1]
	if userID, err := strconv.ParseInt(string(key), 10, 64); err != nil {
		return fmt.Errorf("nsq message corrupted, parse userID error %w", err)
	} else {
		value, ok := s.wsServer.UserConns.Load(userID)
		if !ok {
			log.Debug("message to offline user: %d", userID)
			return nil
		}
		payload := buffer.Bytes()
		if config.Config.Gateway.UseJsonMsg {
			payload, err = protoMessageToJsonMessage(payload)
			if err != nil {
				return fmt.Errorf("can't transfer protobuf message to json: %w", err)
			}
		}
		log.Debug("message to user %d, content-length: %d", userID, len(payload))
		channel := value.(*websocket.Channel)
		if err := channel.Push(payload); err != nil {
			return fmt.Errorf("push message to channel %s error %w", channel.ID(), err)
		}
		return nil
	}
}

func protoMessageToJsonMessage(payload []byte) ([]byte, error) {
	message := &pb.BaseMsg{}
	if err := proto.Unmarshal(payload, message); err != nil {
		return nil, err
	}
	return json.Marshal(message)
}
