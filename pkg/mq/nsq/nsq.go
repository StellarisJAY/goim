package nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
)

type Producer struct {
	pd *nsq.Producer
}

type Consumer struct {
	c *nsq.Consumer
}

func NewProducer() *Producer {
	nsqdAddr := config.Config.Nsq.NsqdAddress
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdAddr, nsqConfig)
	if err != nil {
		panic(fmt.Errorf("create nsq producer error, %w", err))
	}
	// 创建时就连接到broker，避免publish时再建立连接
	if err = producer.Ping(); err != nil {
		panic(fmt.Errorf("can't connect to producer at %s, connect error: %w", nsqdAddr, err))
	}
	return &Producer{pd: producer}
}

func NewConsumer(topic string, channel string, handlers ...nsq.HandlerFunc) *Consumer {
	nsqConfig := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		panic(fmt.Errorf("create nsq consumer error: %w", err))
	}
	// connect之后无法addHandler
	for _, handler := range handlers {
		consumer.AddHandler(handler)
	}
	return &Consumer{c: consumer}
}

func (c *Consumer) Connect() {
	lookupAddrs := config.Config.Nsq.LookupAddresses
	if err := c.c.ConnectToNSQLookupds(lookupAddrs); err != nil {
		panic(fmt.Errorf("connect to nsq lookupd error: %w", err))
	}
}

func (pr *Producer) PushMessage(topic string, key string, value []byte) error {
	return pr.pd.Publish(topic, value)
}

func (pr *Producer) Publish(topic string, body []byte) error {
	return pr.pd.Publish(topic, body)
}

// PublishAsync 异步publish，返回结果channel
func (pr *Producer) PublishAsync(topic string, body []byte) (doneChan chan *nsq.ProducerTransaction, err error) {
	done := make(chan *nsq.ProducerTransaction)
	err = pr.pd.PublishAsync(topic, body, done)
	if err != nil {
		close(done)
		return nil, err
	}
	return doneChan, nil
}

// PublishCallback 异步publish，发布完成后callback
func (pr *Producer) PublishCallback(topic string, body []byte, callback func(*nsq.ProducerTransaction)) error {
	done := make(chan *nsq.ProducerTransaction)
	err := pr.pd.PublishAsync(topic, body, done)
	if err != nil {
		close(done)
		return err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Warn("unexpected error in nsq publish callback %v", err)
			}
		}()
		transaction := <-done
		callback(transaction)
	}()
	return nil
}
