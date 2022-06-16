package sicore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter_Reset(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	wr := newWriter(buf)
	encAddr := fmt.Sprintf("%p\n", wr.enc)
	wr.Reset(nil)
	encAddrAfterReset := fmt.Sprintf("%p\n", wr.enc)

	assert.EqualValues(t, encAddr, encAddrAfterReset)
}

func TestWriter_GetWriter_Reset(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))

	wr := GetWriter(buf, SetJsonEncoder())
	encAddr := fmt.Sprintf("%p\n", wr.enc)

	PutWriter(wr)
	assert.Nil(t, wr.enc)

	wr = GetWriter(buf, SetJsonEncoder())
	encAddrAfterReset := fmt.Sprintf("%p\n", wr.enc)
	assert.NotEqualValues(t, encAddr, encAddrAfterReset)

	PutWriter(wr)
	assert.Nil(t, wr.enc)

	wr = GetWriter(buf)
	if _, ok := wr.enc.(*DefaultEncoder); !ok {
		t.FailNow()
	}
}
