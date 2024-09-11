package sikafka

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
)

// RetryableSyncProducer returns SyncProducer with retry configured.
func RetryableSyncProducer(brokers []string, version, topic string) (*SyncProducer, error) {
	parsedVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, errors.New("sikafka: cannot parse kafka version: " + err.Error())
	}

	config := sarama.NewConfig()
	// config.ClientID = ""
	config.Version = parsedVersion
	config.Producer.Idempotent = true
	config.Producer.Retry.Max = 10
	config.Producer.Retry.BackoffFunc = func(retries int, maxRetries int) time.Duration {
		// linear
		// return time.Duration(retries) * 250 * time.Millisecond

		// exponential
		v := (1 << retries) * 250 * time.Millisecond
		if v > 10000*time.Millisecond {
			v = 10000 * time.Millisecond
		}
		return v
	}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	config.Net.DialTimeout = 10 * time.Second
	config.Net.MaxOpenRequests = 1
	config.Metadata.Retry.Max = 10
	config.Metadata.Retry.BackoffFunc = func(retries int, maxRetries int) time.Duration {
		// linear
		// return time.Duration(retries) * 250 * time.Millisecond

		// exponential
		v := (1 << retries) * 250 * time.Millisecond
		if v > 10000*time.Millisecond {
			v = 10000 * time.Millisecond
		}
		return v
	}
	config.Metadata.RefreshFrequency = 5 * time.Minute

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, errors.New("sikafka: failed to create sync producer: " + err.Error())
	}

	return NewSyncProducer(producer, topic,
		WithSyncProducerOptionRetyMax(2)), nil
}

// RetryableAsyncProducer returns SyncProducer with retry configured.
func RetryableAsyncProducer(brokers []string, version, topic string) (*AsyncProducer, error) {
	parsedVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, errors.New("sikafka: cannot parse kafka version: " + err.Error())
	}

	config := sarama.NewConfig()
	// config.ClientID = ""
	config.Version = parsedVersion
	config.Producer.Idempotent = true
	config.Producer.Retry.Max = 10
	config.Producer.Retry.BackoffFunc = func(retries int, maxRetries int) time.Duration {
		// linear
		// return time.Duration(retries) * 250 * time.Millisecond

		// exponential
		v := (1 << retries) * 250 * time.Millisecond
		if v > 10000*time.Millisecond {
			v = 10000 * time.Millisecond
		}
		return v
	}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	config.Net.DialTimeout = 10 * time.Second
	config.Net.MaxOpenRequests = 1

	config.ChannelBufferSize = 1024
	config.Producer.Timeout = 3000 * time.Millisecond
	config.Producer.Flush.Frequency = 3 * time.Second
	config.Producer.Flush.Messages = 512
	config.Producer.Flush.MaxMessages = 1024

	config.Metadata.Retry.Max = 10
	config.Metadata.Retry.BackoffFunc = func(retries int, maxRetries int) time.Duration {
		// linear
		// return time.Duration(retries) * 250 * time.Millisecond

		// exponential
		v := (1 << retries) * 250 * time.Millisecond
		if v > 10000*time.Millisecond {
			v = 10000 * time.Millisecond
		}
		return v
	}
	config.Metadata.RefreshFrequency = 5 * time.Minute

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errors.New("sikafka: failed to create async producer: " + err.Error())
	}

	return NewAsyncProducer(producer, topic), nil
}

// Deprecated
func AsyncProducerWithConfig(config *sarama.Config, brokers []string) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// Deprecated
func SyncProducerWithConfig(config *sarama.Config, brokers []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// Deprecated
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

// Deprecated
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
