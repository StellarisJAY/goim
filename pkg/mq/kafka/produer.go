package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
)

// Producer 消息生产者
type Producer struct {
	sarama.SyncProducer
	Topic  string
	Addr   []string
	config *sarama.Config
}

func NewProducer(addr []string, topic string) (*Producer, error) {
	p := &Producer{Topic: topic, Addr: addr}
	p.config = sarama.NewConfig()
	p.config.Producer.Return.Successes = true
	p.config.Producer.Return.Errors = true
	// 为了避免从 Producer 到 Broker 的消息丢失，使用 ALL ACK
	p.config.Producer.RequiredAcks = sarama.WaitForAll
	// 使用随机的分区策略
	p.config.Producer.Partitioner = sarama.NewRandomPartitioner
	syncP, err := sarama.NewSyncProducer(p.Addr, p.config)
	if err != nil {
		return nil, err
	}
	p.SyncProducer = syncP
	return p, nil
}

func (p *Producer) PushMessage(topic string, key string, value []byte) error {
	if len(key) == 0 || value == nil || len(value) == 0 {
		return errors.New("key or value can't be empty")
	}
	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := p.SyncProducer.SendMessage(message)
	return err
}
