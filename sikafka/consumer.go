package sikafka

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func ConsumerGroupWithConfig(config *sarama.Config, brokers []string, group string) (sarama.ConsumerGroup, error) {

	client, err := sarama.NewConsumerGroup(brokers, group, config)

	return client, err
}

func DefaultConsumerGroup(brokers []string, group string, version string, assignor string, oldest bool) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	v, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, err
	}
	config.Version = v

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		return nil, errors.New("invalid assignor " + assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	config.Consumer.Group.Session.Timeout = time.Duration(60) * time.Second
	config.Consumer.Group.Heartbeat.Interval = time.Duration(60/3) * time.Second
	config.Consumer.Group.Rebalance.Retry.Max = 10
	config.Consumer.Group.Rebalance.Retry.Backoff = 6 * time.Second

	config.Metadata.Retry.Max = 10
	config.Metadata.Retry.BackoffFunc = func(retries int, maxRetries int) time.Duration {
		v := (1 << retries) * 250 * time.Millisecond
		if v > 10000*time.Millisecond {
			v = 10000 * time.Millisecond
		}
		return v
	}

	// TODO: consumer starts rebalancing right after started, and it freezes
	config.Metadata.RefreshFrequency = 5 * time.Minute

	client, err := sarama.NewConsumerGroup(brokers, group, config)

	return client, err
}

type Consumer interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	MakeReady()
	WaitReady() <-chan bool
	CloseReady()
}

type ConsumerGroup struct {
	sarama.ConsumerGroup
	consumer Consumer
	topics   []string

	isPaused bool
	cancel   context.CancelFunc
}

func NewConsumerGroup(cg sarama.ConsumerGroup, consumer Consumer, topics []string) *ConsumerGroup {
	return &ConsumerGroup{
		ConsumerGroup: cg,
		consumer:      consumer,
		topics:        topics,
		isPaused:      false,
	}
}

func (cg *ConsumerGroup) Toggle() {
	cg.toggleConsumptionFlow()
}

func (cg *ConsumerGroup) StartWith(loaded chan bool) error {
	var consumerCtx context.Context
	consumerCtx, cg.cancel = context.WithCancel(context.Background())

	stopCh := make(chan bool)

	var consumerErr error = nil
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		attempts := 0
		retryMax := 5
		for {
			attempts++

			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			log.Println("consumer group: trying to consume...")
			if err := cg.Consume(consumerCtx, cg.topics, cg.consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) || attempts >= retryMax {
					log.Println("consumer error: " + err.Error())
					consumerErr = err
					cg.cancel()
					close(stopCh)
					return
				}
				if consumerCtx.Err() != nil {
					consumerErr = consumerCtx.Err()
					cg.cancel()
					close(stopCh)
					return
				}

				log.Println("consumer error: retrying: " + err.Error())
				time.Sleep(time.Second * 3)
				continue
				// log.Panicf("Error from consumer: %v\n", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if consumerCtx.Err() != nil {
				consumerErr = consumerCtx.Err()
				cg.cancel()
				close(stopCh)
				return
			}

			cg.consumer.MakeReady()
			attempts = 0
		}
	}(&wg)

	select {
	case <-stopCh:
	case <-cg.consumer.WaitReady(): // Await till the consumer has been set up
		log.Println("sarama: consumer group is up and running")
	}
	loaded <- true

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-consumerCtx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			cg.toggleConsumptionFlow()
		}
	}
	cg.cancel()
	wg.Wait()
	if err := cg.Close(); err != nil {
		log.Println("consumer error: failed to close client: " + err.Error())
		if consumerErr == nil {
			consumerErr = err
		}
	}

	return consumerErr
}

func (cg *ConsumerGroup) Start() error {
	consumerLoaded := make(chan bool, 1)
	defer close(consumerLoaded)
	return cg.StartWith(consumerLoaded)
}

func (cg *ConsumerGroup) Finish() error {
	if cg.cancel == nil {
		return errors.New("ConsumerGroup has not been started")
	}
	cg.cancel()
	return nil
}

func (cg *ConsumerGroup) Stop() error {
	return cg.Finish()
}

func (cg *ConsumerGroup) toggleConsumptionFlow() {
	if cg.isPaused {
		cg.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		cg.PauseAll()
		log.Println("Pausing consumption")
	}

	cg.isPaused = !cg.isPaused
}

// CgConsumer represents a Sarama consumer group consumer
type CgConsumer struct {
	ready      chan bool
	msgHandler MessageHandler
}

func NewCgConsumer(msgHandler MessageHandler) *CgConsumer {
	c := &CgConsumer{}
	c.MakeReady()
	c.msgHandler = msgHandler

	return c
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *CgConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	consumer.CloseReady()
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *CgConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *CgConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("message channel was closed")
				return nil
			}
			// log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			if err := consumer.msgHandler.Handle(message); err == nil {
				session.MarkMessage(message, "")
			}

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *CgConsumer) MakeReady() {
	c.ready = make(chan bool)
}

func (c *CgConsumer) WaitReady() <-chan bool {
	return c.ready
}

func (c *CgConsumer) CloseReady() {
	close(c.ready)
}

type MessageHandler interface {
	Handle(msg *sarama.ConsumerMessage) error
}
