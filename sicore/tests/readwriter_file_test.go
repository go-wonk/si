package sicore_test

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/siutils"
	"github.com/stretchr/testify/assert"
)

func testCreateFileToRead(fileName, data string) error {
	fr, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fr.Close()

	_, err = fr.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func TestReader_File_Read(t *testing.T) {
	fileName := "./data/TestReader_File_Read.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	expected := testDataFile[:10]

	byt := make([]byte, 10)
	n, err := s.Read(byt)
	siutils.AssertNilFail(t, err)

	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)

	fileName2 := "./data/TestReader_Read_2.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName2, testDataFile2))

	f2, err := os.OpenFile(fileName2, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f2.Close()

	s.Reset(f2, sicore.SetDefaultEOFChecker())

	expected = testDataFile2[:10]

	// byt := make([]byte, 10)
	n, err = s.Read(byt)
	siutils.AssertNilFail(t, err)

	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_File_ReadAll(t *testing.T) {
	fileName := "./data/TestReader_File_ReadAll.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	b, err := s.ReadAll()
	siutils.AssertNilFail(t, err)

	str := strings.ReplaceAll(string(b), "\r\n", "\n")
	assert.Equal(t, testDataFile, str)
}

func TestReader_File_ReadSmall(t *testing.T) {
	fileName := "./data/TestReader_File_ReadSmall.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)
	b := make([]byte, 1)
	n, err := s.Read(b)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "{", string(b))
	assert.Equal(t, 1, n)
}

func TestReader_File_ReadZeroCase1(t *testing.T) {
	fileName := "./data/TestReader_File_ReadZeroCase1.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	var b []byte
	n, err := s.Read(b)
	siutils.AssertNilFail(t, err)

	assert.Equal(t, 0, n)
}

func TestReader_File_ReadZeroCase2(t *testing.T) {
	fileName := "./data/TestReader_File_ReadZeroCase2.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetReader(f)
	defer sicore.PutReader(s)

	b := make([]byte, 0, len(testDataFile))
	n, err := s.Read(b)
	siutils.AssertNilFail(t, err)
	assert.Equal(t, 0, n)

}

func TestReader_File_Decode(t *testing.T) {
	fileName := "./data/TestReader_File_Decode.txt"
	siutils.AssertNilFail(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	r := sicore.GetReader(f, sicore.SetJsonDecoder())
	defer sicore.PutReader(r)

	var p Person
	siutils.AssertNilFail(t, r.Decode(&p))
	assert.EqualValues(t, Person{Name: "wonk", Age: 20, Email: "wonk@wonk.org"}, p)
}

func TestWriter_File_Write(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_Write.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"
	n, err := s.Write([]byte(expected))
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())

	assert.EqualValues(t, len(expected), n)
}

func TestWriter_File_WriteMany(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_WriteMany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetWriter(f)
	defer sicore.PutWriter(s)
	line := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	line += "\n"
	expected := bytes.Repeat([]byte(line), 1000)
	n, err := s.Write(expected)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())
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

func TestWriter_File_EncodeDefaultEncoderByte(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeDefaultEncoderByte.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetWriter(f, sicore.SetDefaultEncoder())
	defer sicore.PutWriter(s)
	byt := []byte(`{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`)

	err = s.Encode(byt)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())

	err = s.Encode(&byt)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())

}
func TestWriter_File_EncodeDefaultEncoderString(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeDefaultEncoderString.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	s := sicore.GetWriter(f, sicore.SetDefaultEncoder())
	defer sicore.PutWriter(s)
	str := `{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`

	err = s.Encode(str)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())

	err = s.Encode(&str)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())
}
func TestWriter_File_WriteAnyStruct(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_WriteAnyStruct.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.GetWriter(f, sicore.SetJsonEncoder())
	defer sicore.PutWriter(s)

	err = s.Encode(p)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())
}

func TestWriter_File_EncodeJsonEncodeStruct(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeJsonEncodeStruct.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	// bufio readwriter wrap around f
	bw := bufio.NewWriter(f)
	s := sicore.GetWriter(bw, sicore.SetJsonEncoder())
	defer sicore.PutWriter(s)

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}
	err = s.Encode(p)
	siutils.AssertNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())
}

func TestWriter_File_EncodeNoEncoderFail(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeNoEncoderFail.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sicore.GetWriter(f)
	err = s.Encode(p)
	siutils.AssertNotNilFail(t, err)
	siutils.AssertNilFail(t, s.Flush())

}
