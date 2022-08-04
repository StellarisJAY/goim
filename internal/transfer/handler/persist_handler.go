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
			ID:        msg.Id,
			Content:   []byte(msg.Content),
			Timestamp: msg.Timestamp,
		}
		if msg.Flag == pb.MessageFlag_Group {
			message.Flag = byte(msg.Flag)
			message.User1 = msg.From
			message.User2 = msg.To
		} else {
			if msg.From > msg.To {
				message.User1 = msg.From
				message.User2 = msg.To
				message.Flag = byte(pb.MessageFlag_From)
			} else {
				message.User1 = msg.To
				message.User2 = msg.From
				message.Flag = byte(pb.MessageFlag_To)
			}
		}
		err := dao.InsertMessage(message)
		if err != nil {
			log.Println("persist message error: ", err)
		}
	})
	return nil
}
