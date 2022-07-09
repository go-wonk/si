package sisql_test

import (
	"sync"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sisql"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

func TestSqlDB_Concurrency_QueryMaps(t *testing.T) {

	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db) // sicore.SqlColumn{Name: "id", Type: sicore.SqlColTypeInt},
	// sicore.SqlColumn{Name: "id2", Type: sicore.SqlColTypeInt},

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			for j := 0; j < 10000; j++ {
				query := `
					select $1 as id, $2 as id2
					union all
					select $1+1 as id, $2+1 as id2
				`

				m := make([]map[string]interface{}, 0)
				_, err := sqldb.QueryMaps(query, &m, i, j)
				if err != nil {
					t.Fail()
					break
				}
				if !assert.EqualValues(t, i, m[0]["id"]) {
					t.Fail()
					break
				}
				if !assert.EqualValues(t, j, m[0]["id2"]) {
					t.Fail()
					break
				}
				if !assert.EqualValues(t, i+1, m[1]["id"]) {
					t.Fail()
					break
				}
				if !assert.EqualValues(t, j+1, m[1]["id2"]) {
					t.Fail()
					break
				}
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
}

func TestSqlDB_Concurrency_QueryIntoMapSlice(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db, sisql.WithType("id", sicore.SqlColTypeInt), sisql.WithType("id2", sicore.SqlColTypeInt))

	// sicore.SqlColumn{Name: "id", Type: sicore.SqlColTypeInt},
	// sicore.SqlColumn{Name: "id2", Type: sicore.SqlColTypeInt},
	// sicore.SqlColumn{Name: "decimal_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "numeric_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "char_arr_", Type: sicore.SqlColTypeUints8},

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			for j := 0; j < 100; j++ {
				query := `
					select $1 as id, $2 as id2,
						null as nil,
						'123'::varchar(255) as str,
						123::integer as int2_,
						123::decimal(24,4) as decimal_,
						123::numeric(24,4) as numeric_,
						123::bigint as bigint_,
						'{"abcde", "lunch"}'::char(5)[] as char_arr_,
						'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
						'0123'::bytea as bytea_,
						to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
					union all
					select $1 as id, $2 as id2,
						null as nil,
						'123'::varchar(255) as str,
						123::integer as int2_,
						123::decimal(24,4) as decimal_,
						123::numeric(24,4) as numeric_,
						123::bigint as bigint_,
						'{"abcde", "lunch"}'::char(5)[] as char_arr_,
						'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
						'0123'::bytea as bytea_,
						to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
					union all
					select $1 as id, $2 as id2,
						null as nil,
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
				_, err := sqldb.QueryMaps(query, &m, i, j)
				assert.Nil(t, err)
				if !assert.EqualValues(t, i, m[0]["id"]) {
					break
				}
				if !assert.EqualValues(t, j, m[0]["id2"]) {
					break
				}
				// mb, _ := json.Marshal(m)
				// fmt.Println(string(mb))
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
	// expected := `[{"bigint_":123,"bytea_":"MDEyMw==","char_arr_":"e2FiY2RlLGx1bmNofQ==","decimal_":123,"int2_":123,"nil":null,"numeric_":123,"str":"123","time_":"2022-01-01T12:12:12Z","varchar_arr_":"e2FiY2RlLGx1bmNofQ=="}]`
	// mb, _ := json.Marshal(m)
	// // fmt.Println(string(mb))
	// assert.Equal(t, expected, string(mb))

}
