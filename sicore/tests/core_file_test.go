package sicore_test

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

func TestReadWriter_ReadAllBytes(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"

	b, err := s.ReadAll()
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

func TestReadWriter_ReadSmall(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)
	b := make([]byte, 1)
	n, err := s.Read(b)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)
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

type Person struct {
	Name           string `json:"name"`
	Age            uint8  `json:"age"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	MarriageStatus string `json:"marriage_status"`
	NumChildren    uint8  `json:"num_children"`
}

func TestReadWriter_WriteAnyBytes(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)
	byt := []byte(`{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`)

	n, err := s.WriteAny(byt)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

	n, err = s.WriteAny(&byt)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

}
func TestReadWriter_WriteAnyString(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.NewReadWriter(f)
	str := `{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`

	n, err := s.WriteAny(str)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

	n, err = s.WriteAny(&str)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

}
func TestReadWriter_WriteAnyStruct(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.NewReadWriterWithEncoder(f, sicore.JsonEncoder(f))
	n, err := s.WriteAny(p)
	siutils.NilFail(t, err)

	assert.Equal(t, 1, n)
}

func TestReadWriter_WriteAnyStructFlush(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	// bufio readwriter wrap around f
	br := bufio.NewReader(f)
	bw := bufio.NewWriter(f)
	brw := bufio.NewReadWriter(br, bw)

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.NewReadWriterWithEncoder(brw, sicore.JsonEncoder(brw))
	n, err := s.WriteAny(p)
	siutils.NilFail(t, err)

	assert.Equal(t, 1, n)
}

func TestReadWriter_WriteAnyStructFail(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.NewReadWriter(f)
	n, err := s.WriteAny(p)
	siutils.NotNilFail(t, err)

	assert.Equal(t, 0, n)

}
