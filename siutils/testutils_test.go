package siutils_test

import (
	"errors"
	"testing"

	"github.com/go-wonk/si/siutils"
)

func TestNilFail(t *testing.T) {
	tt := &testing.T{}

	siutils.NilFail(tt, nil)
}

func TestNotNilFail(t *testing.T) {
	tt := &testing.T{}

	siutils.NotNilFail(tt, errors.New("error"))
}

func TestNilFailB(t *testing.T) {
	tt := &testing.B{}

	siutils.NilFailB(tt, nil)
}

func TestNotNilFailB(t *testing.T) {
	tt := &testing.B{}

	siutils.NotNilFailB(tt, errors.New("error"))
}
