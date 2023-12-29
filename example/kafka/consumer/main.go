package main

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-wonk/si/v2/sikafka"
)

type testMessageHandler struct{}

func (h *testMessageHandler) Handle(message *sarama.ConsumerMessage) error {
	fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
	return nil
}
func main() {

	defClient, err := sikafka.DefaultConsumerGroup([]string{"testkafkahost:9092"}, "tp-test-grp1", "3.1.0", "range", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	consumer := sikafka.NewCgConsumer(&testMessageHandler{})
	cg := sikafka.NewConsumerGroup(defClient, consumer, []string{"tp-test-15"})
	cg.Start()
}
