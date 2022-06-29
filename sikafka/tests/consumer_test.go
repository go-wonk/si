package sikafka_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-wonk/si/sikafka"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

var (
	consumedMessage string
)

type messageHandler struct{}

func (h *messageHandler) Handle(message *sarama.ConsumerMessage) {
	log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
	consumedMessage = string(message.Value)
}

func produce(key, value string) error {
	producer, err := sikafka.DefaultSyncProducer([]string{"testkafkahost:9092"})
	if err != nil {
		return err
	}
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-consumer-test")
	p, o, err := sp.Produce([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	fmt.Println(p, o)

	return nil
}

func TestConsumerGroup(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	defClient, err := sikafka.DefaultConsumerGroup([]string{"testkafkahost:9092"}, "tp-consumer-test-grp1", "3.1.0", "range", true)
	siutils.AssertNilFail(t, err)

	consumer := sikafka.NewCgConsumer(&messageHandler{})
	cg := sikafka.NewConsumerGroup(defClient, consumer, []string{"tp-consumer-test"})
	go func() {
		cg.Start()
	}()

	err = produce("test-key", "this-is-a-test-message")
	if !assert.Nil(t, err) {
		cg.Finish()
		t.FailNow()
	}

	time.Sleep(1 * time.Second)
	assert.EqualValues(t, "this-is-a-test-message", consumedMessage)
	siutils.AssertNilFail(t, cg.Finish())
}
