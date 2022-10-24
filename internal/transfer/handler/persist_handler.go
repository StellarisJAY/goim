package handler

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/panjf2000/ants/v2"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"runtime"
)

var persistWorkerPool *ants.Pool

func init() {
	wp, err := ants.NewPool(runtime.NumCPU() * 2)
	if err != nil {
		panic(fmt.Errorf("error occured when creating ants worker pool %w", err))
	}
	persistWorkerPool = wp
}

var PersistMessageHandler = func(message *sarama.ConsumerMessage) error {
	value := message.Value
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	_ = persistWorkerPool.Submit(func() {
		persistTask(msg)
	})
	return nil
}

var NsqPersistHandler = func(message *nsq.Message) error {
	value := message.Body
	msg := new(pb.BaseMsg)
	if err := proto.Unmarshal(value, msg); err != nil {
		return err
	}
	_ = persistWorkerPool.Submit(func() {
		persistTask(msg)
	})
	return nil
}

func persistTask(msg *pb.BaseMsg) {
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
}
