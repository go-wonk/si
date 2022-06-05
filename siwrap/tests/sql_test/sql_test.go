package sql_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrap"
	"github.com/stretchr/testify/assert"
)

func TestSqlDB_QueryRow(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db)

	query := `
		select 
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
	`

	var tim sql.NullTime
	// tl := Table{}
	// tl := TableList{}
	row := sqldb.QueryRow(query)
	err := row.Scan(&tim)
	siutils.AssertNilFail(t, err)

	expected := `2022-01-01 12:12:12 +0000 UTC`
	assert.Equal(t, expected, tim.Time.String())
}

func TestSqlDBQueryStructsSimple(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db)

	query := `
		select * from student limit 10
	`

	// tl := Table{}
	var tl []Student
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)
	fmt.Println(tl)

	// _, err = sqldb.QueryStructs(query, &tl)
	// siutils.AssertNilFail(t, err)
	// fmt.Println(tl)

	// expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"}]`
	// assert.Equal(t, expected, tl.String())
}

type Sample struct {
	Nil  string `json:"nil"`
	Int2 int    `json:"int2_"`
	Int3 *int   `json:"int3_"`
}

func (s *Sample) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func TestSqlDBQueryStructsNil(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db) // sicore.SqlColumn{Name: "decimal_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "numeric_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "char_arr_", Type: sicore.SqlColTypeUints8},

	query := `
		select null as nil, 
			123::integer as int2_,
			234::integer as int3_
		union all
		select null as nil, 
			99123::integer as int2_,
			99234::integer as int3_
	`

	// tl := Table{}
	tl := []Sample{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)
	fmt.Println("tl:", tl[0].String(), tl[1].String())

	// expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"}]`
	// assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_Struct(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db) // sicore.SqlColumn{Name: "decimal_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "numeric_", Type: sicore.SqlColTypeFloat64},
	// sicore.SqlColumn{Name: "char_arr_", Type: sicore.SqlColTypeUints8},

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
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_Slice(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db,
		sicore.SqlColumn{Name: "decimal_", Type: sicore.SqlColTypeFloat64},
		sicore.SqlColumn{Name: "numeric_", Type: sicore.SqlColTypeFloat64},
		sicore.SqlColumn{Name: "char_arr_", Type: sicore.SqlColTypeUints8},
	)

	query := `
		select null as nil, 
			'123'::varchar(255) as str, 
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123456789123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		union all
		select null as nil, 
			'234'::varchar(255) as str, 
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123456789123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121213', 'YYYYMMDDHH24MISS') as time_
	`

	// tl := Table{}
	tl := TableList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:13Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_SliceUseSqlNullType(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db)

	query := `
		select null as nil, 
			'123'::varchar(255) as str, 
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123456789123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123absdfwefasdf'::bytea as bytea_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		union all
		select null as nil, 
			'234'::varchar(255) as str, 
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123456789123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123fewfeasdfzxcv123'::bytea as bytea_,
			to_timestamp('20220101121213', 'YYYYMMDDHH24MISS') as time_
	`

	// tl := Table{}
	tl := TableList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123absdfwefasdf","time_":"2022-01-01T12:12:12Z"},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123fewfeasdfzxcv123","time_":"2022-01-01T12:12:13Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoMap(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db,
		sicore.SqlColumn{Name: "decimal_", Type: sicore.SqlColTypeFloat64},
		sicore.SqlColumn{Name: "numeric_", Type: sicore.SqlColTypeFloat64},
		sicore.SqlColumn{Name: "char_arr_", Type: sicore.SqlColTypeUints8},
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
	_, err := sqldb.QueryMaps(query, &m)
	siutils.AssertNilFail(t, err)

	expected := `[{"bigint_":123,"bytea_":"0123","char_arr_":"e2FiY2RlLGx1bmNofQ==","decimal_":123,"int2_":123,"nil":null,"numeric_":123,"str":"123","time_":"2022-01-01T12:12:12Z","varchar_arr_":"e2FiY2RlLGx1bmNofQ=="}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}

func TestSqlDB_QueryIntoMap_Bool(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db)

	query := `
		select null as nil,
			null as true_1, '1' as true_2, 'Y' as true_3,
			0 as false_1, '0' as false_2, 'N' as false_3
		union all
		select 'abcdef' as nil,
			1 as true_1, '1' as true_2, 'Y' as true_3,
			0 as false_1, '0' as false_2, 'N' as false_3
	`

	m := make([]map[string]interface{}, 0)
	_, err := sqldb.QueryMaps(query, &m)
	siutils.AssertNilFail(t, err)

	expected := `[{"false_1":0,"false_2":"0","false_3":"N","nil":null,"true_1":null,"true_2":"1","true_3":"Y"},{"false_1":0,"false_2":"0","false_3":"N","nil":"abcdef","true_1":1,"true_2":"1","true_3":"Y"}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}

func TestSqlDB_QueryIntoMap_Bool_WithSqlColumn(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := siwrap.NewSqlDB(db,
		sicore.SqlColumn{Name: "true_1", Type: sicore.SqlColTypeBool},
		sicore.SqlColumn{Name: "true_2", Type: sicore.SqlColTypeBool},
		sicore.SqlColumn{Name: "false_1", Type: sicore.SqlColTypeBool},
		sicore.SqlColumn{Name: "false_2", Type: sicore.SqlColTypeBool},
	)

	type BoolTest struct {
		True_1  bool `json:"true_1"`
		True_2  bool `json:"true_2"`
		False_1 bool `json:"false_1"`
		False_2 bool `json:"false_2"`
	}
	query := `
		select null as nil,
			null as true_1, '1' as true_2, 
			0 as false_1, '0' as false_2
		union all
		select null as nil,
			1 as true_1, '1' as true_2,
			0 as false_1, '0' as false_2
	`

	m := make([]map[string]interface{}, 0)
	_, err := sqldb.QueryMaps(query, &m)
	siutils.AssertNilFail(t, err)

	expected := `[{"false_1":false,"false_2":false,"nil":null,"true_1":null,"true_2":true},{"false_1":false,"false_2":false,"nil":null,"true_1":true,"true_2":true}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

	bt := make([]BoolTest, 0)
	err = siutils.DecodeAny(m, &bt)
	siutils.AssertNilFail(t, err)

}
