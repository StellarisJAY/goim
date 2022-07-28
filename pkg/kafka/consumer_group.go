package kafka

import (
	"context"
	"github.com/Shopify/sarama"
)

// ConsumerGroup 消费者组，一条消息只会被相同消费者组的一个消费者消费，每个消费者负责部分partition范围的消息
type ConsumerGroup struct {
	sarama.ConsumerGroup
	Topics  []string
	Addr    []string
	GroupID string
}

func NewConsumerGroup(groupID string, addr []string, topics []string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	// 暂时配置从最新的offset消费
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	cg, err := sarama.NewConsumerGroup(addr, groupID, config)
	if err != nil {
		return nil, err
	}
	return &ConsumerGroup{
		ConsumerGroup: cg,
		Topics:        topics,
		Addr:          addr,
		GroupID:       groupID,
	}, nil
}

func (cg *ConsumerGroup) Start(ctx context.Context, handle MessageHandleFunc) {
	// 因为新的消费者加入导致的re-balance会使 Consume 退出，所以需要循环
	for {
		err := cg.Consume(ctx, cg.Topics, &baseConsumerGroupHandler{handle: handle})
		if err != nil {
			panic(err.Error())
		}
	}
}
