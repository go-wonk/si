package siutils_test

import (
	"testing"

	"github.com/go-wonk/si/siutils"
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
