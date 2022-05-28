package sicore

import (
	"testing"

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
