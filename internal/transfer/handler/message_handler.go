package handler

import (
	"github.com/Shopify/sarama"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"log"
)

var MessageTransferHandler = func(message *sarama.ConsumerMessage) error {
	value := message.Value
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	log.Printf("transfer received message, from: %d, to: %d, content: %s, flag: %d", msg.From, msg.To, msg.Content, msg.Flag)
	return nil
}
