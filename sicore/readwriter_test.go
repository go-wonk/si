package sicore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-wonk/si/siutils"
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
	err := wr.EncodeFlush([]byte("test message"))
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "test message", buf.String())

	PutWriter(wr)

	var buf2 bytes.Buffer
	wr2 := GetWriter(&buf2)
	err = wr2.EncodeFlush("test message 2")
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "test message 2", buf2.String())
	PutWriter(wr2)
}

func TestReader_Reset(t *testing.T) {

	buf := bytes.NewBuffer([]byte("test message"))
	rd := GetReader(buf)

	var res []byte
	err := rd.Decode(&res)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))
	PutReader(rd)

	buf2 := bytes.NewBuffer([]byte("test message2"))
	rd2 := GetReader(buf2)

	var res2 []byte
	err = rd2.Decode(&res2)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res2))
	PutReader(rd2)

}
