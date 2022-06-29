package sikafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func DefaultAsyncProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal

	// config.Producer.Timeout = 2000 * time.Millisecond
	config.Producer.Flush.Frequency = 1000 * time.Millisecond
	config.Producer.Flush.Messages = 1000
	config.Producer.Flush.MaxMessages = 4096
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
func DefaultSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 50
	config.ChannelBufferSize = 4096
	config.Producer.Timeout = 3000 * time.Millisecond
	// config.Producer.Flush.Frequency = 250 * time.Millisecond
	// config.Producer.Flush.Messages = 1000

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

type SyncProducer struct {
	sarama.SyncProducer
	topic string
}

func NewSyncProducer(producer sarama.SyncProducer, topic string) *SyncProducer {
	return &SyncProducer{producer, topic}
}

func (sp *SyncProducer) Produce(key []byte, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: sp.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err = sp.SendMessage(msg)
	return
}
func (sp *SyncProducer) ProduceWithTopic(topic string, key []byte, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err = sp.SendMessage(msg)
	return
}

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
