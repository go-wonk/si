package sql_test

import (
	"encoding/json"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/mitchellh/mapstructure"
)

var (
	bmap = map[string]interface{}{
		"nil":   "",
		"int2_": 123,
	}
)

func BenchmarkDecode_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		table := Table{}
		byt, _ := json.Marshal(bmap)
		err := json.Unmarshal(byt, &table)
		siutils.NilFailB(b, err)
	}
}
func BenchmarkDecode_Mapstructure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		table := Table{}
		err := mapstructure.Decode(bmap, &table)
		siutils.NilFailB(b, err)
	}
}
