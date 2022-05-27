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

func TestReader_Read(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	expected := `{"name":"w`

	byt := make([]byte, 10)
	n, err := s.Read(byt)
	siutils.NilFail(t, err)

	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_ReadAll(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"

	b, err := s.ReadAll()
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	str := strings.ReplaceAll(string(b), "\r\n", "\n")
	assert.Equal(t, expected, str)
}

func TestReadWriter_ReadSmall(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)
	b := make([]byte, 1)
	n, err := s.Read(b)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)
}

func TestReader_ReadZeroCase1(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)
	var b []byte
	n, err := s.Read(b)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 0, n)
}

func TestReader_ReadZeroCase2(t *testing.T) {
	f, _ := os.OpenFile("./data/readonly.txt", os.O_RDONLY, 0644)
	defer f.Close()

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)
	b := make([]byte, 0, len(expected))
	n, err := s.Read(b)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 0, n)
}

func TestWriter_Write(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_Write.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.NilFail(t, err)

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"
	n, err := s.Write([]byte(expected))
	siutils.NilFail(t, err)

	assert.EqualValues(t, len(expected), n)
}

func TestWriter_WriteMany(t *testing.T) {
	f, _ := os.OpenFile("./data/TestWriter_WriteMany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)
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

func TestWriter_WriteAnyBytes(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)
	byt := []byte(`{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`)

	n, err := s.WriteAny(byt)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

	n, err = s.WriteAny(&byt)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

}
func TestWriter_WriteAnyString(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)
	str := `{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`

	n, err := s.WriteAny(str)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

	n, err = s.WriteAny(&str)
	siutils.NilFail(t, err)
	assert.Equal(t, 1, n)

}
func TestWriter_WriteAnyStruct(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.GetWriterWithEncoder(f, sicore.JsonEncoder(f))
	defer sicore.PutWriter(s)

	n, err := s.WriteAny(p)
	siutils.NilFail(t, err)

	assert.Equal(t, 1, n)
}

func TestWriter_WriteAnyStructFlush(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	// bufio readwriter wrap around f
	bw := bufio.NewWriter(f)
	s := sicore.GetWriterSizeWithEncoder(bw, 1024, sicore.JsonEncoder(bw))
	defer sicore.PutWriter(s)

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}
	n, err := s.WriteAny(p)
	siutils.NilFail(t, err)

	assert.Equal(t, 1, n)
}

func TestWriter_WriteAnyStructFail(t *testing.T) {
	f, _ := os.OpenFile("./data/writeany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.GetWriter(f)
	n, err := s.WriteAny(p)
	siutils.NotNilFail(t, err)

	assert.Equal(t, 0, n)

}
