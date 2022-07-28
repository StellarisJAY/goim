package kafka

import (
	"github.com/Shopify/sarama"
)

// MessageHandleFunc 消息处理器
type MessageHandleFunc func(message *sarama.ConsumerMessage) error

// 基础的消费组处理器，对业务层省略 Setup 和 Cleanup
type baseConsumerGroupHandler struct {
	handle MessageHandleFunc
}

func (b *baseConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *baseConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// 消费者组 re-balance 会触发Cleanup
	return nil
}

func (b *baseConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 接受并处理消息
	for message := range claim.Messages() {
		if err := b.handle(message); err != nil {
			return err
		}
	}
	return nil
}
