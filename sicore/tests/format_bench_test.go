package sicore_test

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%07d", 100)
	}
}

func pad(input []byte, c byte, count int) []byte {
	if len(input) >= count {
		return input
	}
	d := count - len(input)
	p := make([]byte, count)
	for i := 0; i < d; i++ {
		p[i] = c
	}
	for i, v := range input {
		p[d+i] = v
	}
	return p
}

func BenchmarkFormat2(b *testing.B) {
	for i := 0; i < b.N; i++ {

		s := []byte(strconv.FormatInt(100, 10))
		_ = pad(s, '0', 7)

	}
}

func BenchmarkFormat3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseInt("0000100", 10, 64)
	}
}
