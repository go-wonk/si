package sicore_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/stretchr/testify/assert"
)

func TestReadWriter_ReadAllBytes(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"

	b, err := s.ReadAllBytes()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	str := strings.ReplaceAll(string(b), "\r\n", "\n")
	assert.Equal(t, expected, str)
}

func TestReadWriter_Read(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`

	s := sicore.NewReadWriter(f)
	b := make([]byte, len(expected))
	n, err := s.Read(b)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, len(expected), n)
}

func TestReadWriter_ReadZeroCase1(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)
	var b []byte
	n, err := s.Read(b)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 0, n)
}

func TestReadWriter_ReadZeroCase2(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`

	s := sicore.NewReadWriter(f)
	b := make([]byte, 0, len(expected))
	n, err := s.Read(b)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 0, n)
}

func TestReadWriter_Write(t *testing.T) {
	f, _ := os.OpenFile("./data/write.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)
	line := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	line += "\n"
	expected := bytes.Repeat([]byte(line), 1000)
	n, err := s.Write(expected)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, len(line)*1000, n)

}
