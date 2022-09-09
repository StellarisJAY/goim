package handler

import (
	"github.com/Shopify/sarama"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/pool"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"runtime"
)

var persistWorkers = pool.NewWorkerPool(runtime.NumCPU() * 2)

var persistChannel = make(chan *model.Message, 1024)

func init() {
	persistWorkers.Submit(func() {
		for msg := range persistChannel {
			size := len(persistChannel)
			messages := make([]*model.Message, size+1)
			messages[0] = msg
			for i := 0; i < size; i++ {
				messages[i+1] = <-persistChannel
			}
			err := dao.InsertMessages(messages)
			if err != nil {
				log.Warn("persist messages error: ", err)
			}
		}
	})
}

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
			log.Warn("persist message error: ", err)
		}
	})
	return nil
}
