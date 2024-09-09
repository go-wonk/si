package sikafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/eapache/go-resiliency/breaker"
)

func AsyncProducerWithConfig(config *sarama.Config, brokers []string) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
func SyncProducerWithConfig(config *sarama.Config, brokers []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

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
	partition, offset, err = sp.produce(msg)
	return
}
func (sp *SyncProducer) ProduceWithTopic(topic string, key []byte, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err = sp.produce(msg)
	return
}

func (sp *SyncProducer) ProduceWithMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	partition, offset, err = sp.produce(msg)
	return
}

func (sp *SyncProducer) produce(message *sarama.ProducerMessage) (int32, int64, error) {
	var err error
	var partition int32
	var offset int64

	attempts := 0
	wait := 100 * time.Millisecond
	for attempts < 8 {
		partition, offset, err = sp.SendMessage(message)
		if err != nil {
			if isRetryableError(err) {
				time.Sleep(wait)
				wait *= 2
				continue
			}
			return partition, offset, err
		}
		return partition, offset, nil
	}
	return partition, offset, err
}

func isRetryableError(err error) bool {
	switch err {
	case sarama.ErrBrokerNotAvailable,
		sarama.ErrLeaderNotAvailable,
		sarama.ErrReplicaNotAvailable,
		sarama.ErrRequestTimedOut,
		sarama.ErrNotEnoughReplicas,
		// sarama.ErrNotEnoughReplicasAfterAppend, // "kafka server: Messages are written to the log, but to fewer in-sync replicas than required"
		// sarama.ErrNetworkException, // "kafka server: The server disconnected before a response was received"
		sarama.ErrOutOfBrokers,
		sarama.ErrOutOfOrderSequenceNumber,
		sarama.ErrNotController,
		sarama.ErrNotLeaderForPartition,
		breaker.ErrBreakerOpen:
		return true
	default:
		return false
	}
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

func (ap *AsyncProducer) ProduceWithMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {

	ap.Input() <- msg

	return 0, 0, nil
}
