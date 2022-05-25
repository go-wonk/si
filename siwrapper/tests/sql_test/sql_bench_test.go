package sql_test

import (
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrapper"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSqlDB_QueryIntoMap(b *testing.B) {
	if onlinetest != "1" {
		b.Skip("skipping online tests")
	}
	siutils.NotNilFailB(b, db)

	for i := 0; i < b.N; i++ {
		sqldb := siwrapper.NewSqlDB(db)
		sqldb.AddSqlColumn(sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
			sicore.SqlColumn{"numeric_", sicore.SqlColTypeFloat64},
			sicore.SqlColumn{"char_arr_", sicore.SqlColTypeUints8},
		)

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		m := make([]map[string]interface{}, 0)
		_, err := sqldb.QueryIntoMap(query, &m)
		siutils.NilFailB(b, err)
	}

}

/*
Benchmark on json and mapstructure

goos: windows
goarch: amd64
pkg: github.com/go-wonk/si/siwrapper/tests/sql_test
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkSqlDB_QueryIntoAny_Struct-8                 100           1683737 ns/op            7236 B/op        151 allocs/op
BenchmarkSqlDB_QueryIntoAny_Struct2-8                100           2903963 ns/op           11857 B/op        252 allocs/op
PASS
*/
func BenchmarkSqlDB_QueryIntoAny_Struct(b *testing.B) {
	if onlinetest != "1" {
		b.Skip("skipping online tests")
	}
	siutils.NotNilFailB(b, db)

	for i := 0; i < b.N; i++ {
		sqldb := siwrapper.NewSqlDB(db)
		sqldb.AddSqlColumn(sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
			sicore.SqlColumn{"numeric_", sicore.SqlColTypeFloat64},
			sicore.SqlColumn{"char_arr_", sicore.SqlColTypeUints8},
		)

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		// tl := Table{}
		tl := TableList{}
		_, err := sqldb.QueryIntoAny(query, &tl)
		siutils.NilFailB(b, err)

		expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:12Z"}]`
		assert.Equal(b, expected, tl.String())
	}
}
