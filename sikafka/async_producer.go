package sikafka

import (
	"github.com/IBM/sarama"
)

type AsyncProducer struct {
	sarama.AsyncProducer
	topic string
}

func NewAsyncProducer(producer sarama.AsyncProducer, topic string) *AsyncProducer {
	return &AsyncProducer{producer, topic}
}

func (ap *AsyncProducer) Produce(key []byte, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: ap.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	ap.Input() <- msg

	return 0, 0, nil
}

func (ap *AsyncProducer) ProduceWithTopic(topic string, key []byte, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	ap.Input() <- msg

	return 0, 0, nil
}

func (ap *AsyncProducer) ProduceWithMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {

	ap.Input() <- msg

	return 0, 0, nil
}
