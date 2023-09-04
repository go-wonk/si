package siutils_test

import (
	"testing"

	"github.com/go-wonk/si/v2/siutils"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAny(t *testing.T) {
	m := map[string]interface{}{
		"name": "wonk",
		"age":  20,
	}
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{}
	err := siutils.DecodeAny(m, &p)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, Person{"wonk", 20}, p)
}

func TestDecodeAnyJsonIter(t *testing.T) {
	m := map[string]interface{}{
		"name": "wonk",
		"age":  20,
	}
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{}

	err := siutils.DecodeAny(m, &p)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, Person{Name: "wonk", Age: 20}, p)
}

func BenchmarkDecodeAny(b *testing.B) {
	m := map[string]interface{}{
		"name": "wonk",
		"age":  20,
		"a1":   1234.1234567,
		"a2":   1234.1234567,
		"a3":   1234.1234567,
		"a4":   1234.1234567,
		"a5":   1234.1234567,
		"a6":   1234.1234567,
		"a7":   1234.1234567,
		"a8":   1234.1234567,
		"a9":   1234.1234567,
		"a10":  1234.1234567,
		"a11":  1234.1234567,
	}
	type Person struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		A1   float64 `json:"a1"`
		A2   float64 `json:"a2"`
		A3   float64 `json:"a3"`
		A4   float64 `json:"a4"`
		A5   float64 `json:"a5"`
		A6   float64 `json:"a6"`
		A7   float64 `json:"a7"`
		A8   float64 `json:"a8"`
		A9   float64 `json:"a9"`
		A10  float64 `json:"a10"`
		A11  float64 `json:"a11"`
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p := Person{}
		err := siutils.DecodeAny(m, &p)
		siutils.AssertNilFailB(b, err)
	}

	// assert.EqualValues(b, Person{Name: "wonk", Age: 20}, p)
}

func BenchmarkDecodeAnyJsonIter(b *testing.B) {
	m := map[string]interface{}{
		"name": "wonk",
		"age":  20,
		"a1":   1234.1234567,
		"a2":   1234.1234567,
		"a3":   1234.1234567,
		"a4":   1234.1234567,
		"a5":   1234.1234567,
		"a6":   1234.1234567,
		"a7":   1234.1234567,
		"a8":   1234.1234567,
		"a9":   1234.1234567,
		"a10":  1234.1234567,
		"a11":  1234.1234567,
	}
	type Person struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		A1   float64 `json:"a1"`
		A2   float64 `json:"a2"`
		A3   float64 `json:"a3"`
		A4   float64 `json:"a4"`
		A5   float64 `json:"a5"`
		A6   float64 `json:"a6"`
		A7   float64 `json:"a7"`
		A8   float64 `json:"a8"`
		A9   float64 `json:"a9"`
		A10  float64 `json:"a10"`
		A11  float64 `json:"a11"`
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p := Person{}
		err := siutils.DecodeAny(m, &p)
		siutils.AssertNilFailB(b, err)
	}
	// assert.EqualValues(t, Person{"wonk", 20}, p)
}
