package sikafka

import (
	"time"

	"github.com/IBM/sarama"
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
