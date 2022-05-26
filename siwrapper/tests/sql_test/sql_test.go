package sql_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrapper"
	"github.com/stretchr/testify/assert"
)

func TestSqlDB_QueryRow(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db)

	query := `
		select 
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
	`

	var tim sql.NullTime
	// tl := Table{}
	// tl := TableList{}
	row := sqldb.QueryRow(query)
	err := row.Scan(&tim)
	siutils.NilFail(t, err)

	expected := `2022-01-01 12:12:12 +0000 UTC`
	assert.Equal(t, expected, tim.Time.String())
}

func TestSqlDB_QueryIntoAny_Struct(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db, sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
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
	siutils.NilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_Slice(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db, sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
		sicore.SqlColumn{"numeric_", sicore.SqlColTypeFloat64},
		sicore.SqlColumn{"char_arr_", sicore.SqlColTypeUints8},
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
	_, err := sqldb.QueryIntoAny(query, &tl)
	siutils.NilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:13Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_SliceUseSqlNullType(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db)

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
	_, err := sqldb.QueryIntoAny(query, &tl)
	siutils.NilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123absdfwefasdf","time_":"2022-01-01T12:12:12Z"},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123fewfeasdfzxcv123","time_":"2022-01-01T12:12:13Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoMap(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db, sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
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
	_, err := sqldb.QueryIntoMapSlice(query, &m)
	siutils.NilFail(t, err)

	expected := `[{"bigint_":123,"bytea_":"0123","char_arr_":"e2FiY2RlLGx1bmNofQ==","decimal_":123,"int2_":123,"nil":null,"numeric_":123,"str":"123","time_":"2022-01-01T12:12:12Z","varchar_arr_":"e2FiY2RlLGx1bmNofQ=="}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}

func TestSqlDB_QueryIntoMap_Bool(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db)

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
	_, err := sqldb.QueryIntoMapSlice(query, &m)
	siutils.NilFail(t, err)

	expected := `[{"false_1":0,"false_2":"0","false_3":"N","nil":null,"true_1":null,"true_2":"1","true_3":"Y"},{"false_1":0,"false_2":"0","false_3":"N","nil":"abcdef","true_1":1,"true_2":"1","true_3":"Y"}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}

func TestSqlDB_QueryIntoMap_Bool_WithSqlColumn(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db, sicore.SqlColumn{"true_1", sicore.SqlColTypeBool},
		sicore.SqlColumn{"true_2", sicore.SqlColTypeBool},
		sicore.SqlColumn{"false_1", sicore.SqlColTypeBool},
		sicore.SqlColumn{"false_2", sicore.SqlColTypeBool},
	)

	type BoolTest struct {
		true_1  bool
		true_2  bool
		false_1 bool
		false_2 bool
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
	_, err := sqldb.QueryIntoMapSlice(query, &m)
	siutils.NilFail(t, err)

	expected := `[{"false_1":false,"false_2":false,"nil":null,"true_1":null,"true_2":true},{"false_1":false,"false_2":false,"nil":null,"true_1":true,"true_2":true}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

	bt := make([]BoolTest, 0)
	err = siutils.DecodeAny(m, &bt)
	siutils.NilFail(t, err)

}
