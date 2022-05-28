package sicore

import (
	"os"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

func Test_grow(t *testing.T) {
	b := make([]byte, 0, 10)
	b = append(b, []byte("asdfㅁ")...)
	assert.Equal(t, 7, len(b))
	assert.Equal(t, 10, cap(b))
	l, err := grow(&b, 100)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 7, l)
	assert.Equal(t, 107, len(b))
	assert.Equal(t, 120, cap(b))
}

func Test_growCap(t *testing.T) {
	b := make([]byte, 0, 10)
	b = append(b, []byte("asdfㅁ")...)
	assert.Equal(t, 7, len(b))
	assert.Equal(t, 10, cap(b))

	err := growCap(&b, 100)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 7, len(b))
	assert.Equal(t, 110, cap(b))
}

func Benchmark_readAll_4096(b *testing.B) {
	f, err := os.OpenFile("./tests/data/readonly.txt", os.O_RDONLY, 0644)
	siutils.NilFailB(b, err)

	for i := 0; i < b.N; i++ {
		_, err := readAll(f, DefaultValidator())
		siutils.NilFailB(b, err)
	}
}

func BenchmarkBytesReadWriter_readAll_1024(b *testing.B) {
	f, err := os.OpenFile("./tests/data/readonly.txt", os.O_RDONLY, 0644)
	siutils.NilFailB(b, err)

	for i := 0; i < b.N; i++ {
		_, err := readAll(f, DefaultValidator())
		siutils.NilFailB(b, err)
	}
}
