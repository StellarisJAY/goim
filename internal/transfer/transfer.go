package transfer

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stellarisJAY/goim/internal/transfer/handler"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/mq/kafka"
	"github.com/stellarisJAY/goim/pkg/mq/nsq"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"net/http"
	"strings"
)

var (
	transferConsumerGroup *kafka.ConsumerGroup
	persistConsumerGroup  *kafka.ConsumerGroup
	nsqTransferConsumer   *nsq.Consumer
	nsqPersistConsumer    *nsq.Consumer
)

func Init() {
	mq := strings.ToLower(config.Config.MessageQueue)
	switch mq {
	case "kafka":
		transferConsumerGroup = kafka.NewConsumerGroup(pb.MessageTransferGroup, config.Config.Kafka.Addrs, []string{pb.MessageTransferTopic})
		persistConsumerGroup = kafka.NewConsumerGroup(pb.MessagePersistGroup, config.Config.Kafka.Addrs, []string{pb.MessageTransferTopic})
	case "nsq":
		nsqTransferConsumer = nsq.NewConsumer(pb.MessageTransferTopic, pb.MessageTransferGroup, handler.NsqMessageHandler)
		nsqPersistConsumer = nsq.NewConsumer(pb.MessageTransferTopic, pb.MessagePersistGroup, handler.NsqPersistHandler)
	default:
		panic(fmt.Errorf("unknown or unsupported message queue %s", mq))
	}
}

func Start() {
	ctx := context.Background()
	mq := strings.ToLower(config.Config.MessageQueue)
	switch mq {
	case "kafka":
		go transferConsumerGroup.Start(ctx, handler.MessageTransferHandler)
		go persistConsumerGroup.Start(ctx, handler.PersistMessageHandler)
	case "nsq":
		nsqTransferConsumer.Connect()
		nsqPersistConsumer.Connect()
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(config.Config.Metrics.PromHttpAddr, nil)
	}()
	<-ctx.Done()
}
