package siutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NilFail(t *testing.T, v any) {
	if !assert.Nil(t, v) {
		t.FailNow()
	}
}
func NotNilFail(t *testing.T, v any) {
	if !assert.NotNil(t, v) {
		t.FailNow()
	}
}

func NilFailB(t *testing.B, v any) {
	if !assert.Nil(t, v) {
		t.FailNow()
	}
}
func NotNilFailB(t *testing.B, v any) {
	if !assert.NotNil(t, v) {
		t.FailNow()
	}
}
