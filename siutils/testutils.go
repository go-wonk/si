package siutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertNilFail(t *testing.T, v any) {
	if !assert.Nil(t, v) {
		t.FailNow()
	}
}
func AssertNotNilFail(t *testing.T, v any) {
	if !assert.NotNil(t, v) {
		t.FailNow()
	}
}

func AssertNilFailB(t *testing.B, v any) {
	if !assert.Nil(t, v) {
		t.FailNow()
	}
}
func AssertNotNilFailB(t *testing.B, v any) {
	if !assert.NotNil(t, v) {
		t.FailNow()
	}
}
