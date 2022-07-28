package handler

import (
	"github.com/Shopify/sarama"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/pool"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"runtime"
)

var persistWorkers = pool.NewWorkerPool(runtime.NumCPU() * 2)

var PersistMessageHandler = func(message *sarama.ConsumerMessage) error {
	value := message.Value
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	persistWorkers.Submit(func() {
		message := &model.Message{
			Content:   []byte(msg.Content),
			Timestamp: msg.Timestamp,
		}
		if byte(msg.Flag) == pb.MessageFlagGroup {
			message.Flag = byte(msg.Flag)
			message.User1 = msg.From
			message.User2 = msg.To
		} else {
			if msg.From > msg.To {
				message.User1 = msg.From
				message.User2 = msg.To
				message.Flag = pb.MessageFlagFrom
			} else {
				message.User1 = msg.To
				message.User2 = msg.From
				message.Flag = pb.MessageFlagTo
			}
		}
		err := dao.InsertMessage(message)
		if err != nil {
			log.Println("persist message error: ", err)
		}
	})
	return nil
}
