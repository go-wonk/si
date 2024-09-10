package sikafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/eapache/go-resiliency/breaker"
)

const defaultRetryMax uint16 = 1

type SyncProducer struct {
	sarama.SyncProducer
	topic    string
	retryMax uint16
}

func NewSyncProducer(producer sarama.SyncProducer, topic string, opts ...SyncProducerOption) *SyncProducer {
	p := &SyncProducer{producer, topic, defaultRetryMax}
	for _, o := range opts {
		if o == nil {
			continue
		}
		_ = o.apply(p)
	}
	return p
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
	for attempts < int(sp.retryMax) {
		attempts++
		partition, offset, err = sp.SendMessage(message)
		if err != nil {
			if isRetryableError(err) && attempts < int(sp.retryMax) {
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
