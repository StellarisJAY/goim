package transfer

import (
	"context"
	"github.com/stellarisJAY/goim/internal/transfer/handler"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/kafka"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

var (
	consumer     *kafka.ConsumerGroup
	persistGroup *kafka.ConsumerGroup
)

func Init() {
	consumer = kafka.NewConsumerGroup(pb.MessageTransferGroup, config.Config.Kafka.Addrs, []string{pb.MessageTransferTopic})
	persistGroup = kafka.NewConsumerGroup(pb.MessagePersistGroup, config.Config.Kafka.Addrs, []string{pb.MessageTransferTopic})
}

func Start() {
	ctx := context.Background()
	go consumer.Start(ctx, handler.MessageTransferHandler)
	go persistGroup.Start(ctx, handler.PersistMessageHandler)
	<-ctx.Done()
}
