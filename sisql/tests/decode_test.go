package sisql_test

import (
	"encoding/json"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
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
		table := testmodels.Table{}
		byt, _ := json.Marshal(bmap)
		err := json.Unmarshal(byt, &table)
		siutils.AssertNilFailB(b, err)
	}
}
func BenchmarkDecode_Mapstructure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		table := testmodels.Table{}
		err := mapstructure.Decode(bmap, &table)
		siutils.AssertNilFailB(b, err)
	}
}
