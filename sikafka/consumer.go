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
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	default:
		return nil, errors.New("invalid assignor " + assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	client, err := sarama.NewConsumerGroup(brokers, group, config)

	return client, err
}

type Consumer interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	MakeReady()
	WaitReady()
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

func (cg *ConsumerGroup) Toggle() {
	cg.toggleConsumptionFlow()
}

func (cg *ConsumerGroup) StartWith(loaded chan bool) error {
	var ctx context.Context
	ctx, cg.cancel = context.WithCancel(context.Background())

	var startErr error = nil
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := cg.Consume(ctx, cg.topics, cg.consumer); err != nil {
				// TODO: handle error
				// log.Panicf("Error from consumer: %v", err)
				for i := 0; i < 12; i++ {
					time.Sleep(time.Second * 1)
					if ctx.Err() != nil {
						startErr = ctx.Err()
						return
					}
				}
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				startErr = ctx.Err()
				return
			}
			// consumer.ready = make(chan bool)
			cg.consumer.MakeReady()
		}
	}(&wg)

	cg.consumer.WaitReady() // Await till the consumer has been set up
	loaded <- true
	log.Println("Sarama consumer up and running")

	// sigusr1 := make(chan os.Signal, 1)
	// signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
			// case <-sigusr1:
			// 	cg.toggleConsumptionFlow()
		}
	}
	cg.cancel()
	wg.Wait()
	if err := cg.Close(); err != nil {
		if startErr == nil {
			startErr = err
		}
		log.Panicf("Error closing client: %v", err)
	}

	return startErr
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

func (c *CgConsumer) WaitReady() {
	<-c.ready
}

func (c *CgConsumer) CloseReady() {
	close(c.ready)
}

type MessageHandler interface {
	Handle(msg *sarama.ConsumerMessage) error
}
