package siutils_test

import (
	"errors"
	"testing"

	"github.com/go-wonk/si/siutils"
)

func TestNilFail(t *testing.T) {
	tt := &testing.T{}

	siutils.AssertNilFail(tt, nil)
}

func TestNotNilFail(t *testing.T) {
	tt := &testing.T{}

	siutils.AssertNotNilFail(tt, errors.New("error"))
}

func TestNilFailB(t *testing.T) {
	tt := &testing.B{}

	siutils.AssertNilFailB(tt, nil)
}

func TestNotNilFailB(t *testing.T) {
	tt := &testing.B{}

	siutils.AssertNotNilFailB(tt, errors.New("error"))
}
