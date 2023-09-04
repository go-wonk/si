package sicore_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/siutils"
	"github.com/stretchr/testify/assert"
)

func TestReader_Buffer_Read(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	r := sicore.GetReader(buf)
	defer sicore.PutReader(r)

	expected := testDataFile[:10]
	byt := make([]byte, 10)
	n, err := r.Read(byt)
	siutils.AssertNilFail(t, err)
	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_Buffer_ReadBufio(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	br := bufio.NewReader(buf)

	r := sicore.GetReader(br)
	defer sicore.PutReader(r)

	expected := testDataFile[:10]
	byt := make([]byte, 10)
	n, err := r.Read(byt)
	siutils.AssertNilFail(t, err)
	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_Buffer_ReadAll(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	r := sicore.GetReader(buf)
	defer sicore.PutReader(r)

	expected := testDataFile

	byt, err := r.ReadAll()
	siutils.AssertNilFail(t, err)
	assert.Equal(t, expected, string(byt))
}
