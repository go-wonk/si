package sicore

import (
	"os"
	"testing"

	"github.com/go-wonk/si/siutils"
)

func testCreateFileToRead(fileName, data string) error {
	fr, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

func Benchmark_readAll(b *testing.B) {
	fileName := "./tests/data/Benchmark_readAll.txt"
	siutils.AssertNilFailB(b, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	siutils.AssertNilFailB(b, err)
	defer f.Close()

	for i := 0; i < b.N; i++ {
		_, err := readAll(f, DefaultEofChecker)
		siutils.AssertNilFailB(b, err)
	}
}
