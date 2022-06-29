package sikafka_test

import (
	"strconv"
	"testing"

	"github.com/go-wonk/si/sikafka"
	"github.com/go-wonk/si/siutils"
)

func BenchmarkSyncProducer_Produce(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	producer, err := sikafka.DefaultSyncProducer([]string{"testkafkahost:9092"})
	siutils.AssertNilFailB(b, err)
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-test-15")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := strconv.FormatInt(int64(i), 10)
		sp.Produce([]byte(key), []byte("asdf"))
		// siutils.AssertNilFailB(b, err)
	}
}

func BenchmarkAsyncProducer_Produce(b *testing.B) {

	producer, err := sikafka.DefaultAsyncProducer([]string{"testkafkahost:9092"})
	siutils.AssertNilFailB(b, err)
	defer producer.AsyncClose()

	sp := sikafka.NewAsyncProducer(producer, "tp-test-15")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := strconv.FormatInt(int64(i), 10)
		sp.Produce([]byte(key), []byte("asdf"))
		// siutils.AssertNilFailB(b, err)
	}
}
