package mq

type MessageProducer interface {
	PushMessage(topic, key string, value []byte) error
}
