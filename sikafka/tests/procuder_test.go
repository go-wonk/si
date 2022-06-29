package sikafka_test

import (
	"fmt"
	"testing"

	"github.com/go-wonk/si/sikafka"
	"github.com/go-wonk/si/siutils"
)

func TestProducer_Produce(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	producer, err := sikafka.DefaultSyncProducer([]string{"kafkahost:9092"})
	siutils.AssertNilFail(t, err)
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-test-15")
	p, o, err := sp.Produce([]byte("10123"), []byte("asdf"))
	siutils.AssertNilFail(t, err)
	fmt.Println(p, o)
}

func TestProducer_ProduceWithTopic(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	producer, err := sikafka.DefaultSyncProducer([]string{"kafkahost:9092"})
	siutils.AssertNilFail(t, err)
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-test-15")
	p, o, err := sp.ProduceWithTopic("tp-test", []byte("10123"), []byte("asdf"))
	siutils.AssertNilFail(t, err)
	fmt.Println(p, o)
}
