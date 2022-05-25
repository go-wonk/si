package sql_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrapper"
	"github.com/stretchr/testify/assert"
)

// func TestSqlDB_Query(t *testing.T) {
// 	if onlinetest != "1" {
// 		t.Skip("skipping online tests")
// 	}
// 	siutils.AssertNotNilAndFilaNow(db, t)

// 	sqldb := siwrapper.NewSqlDB(db)

// 	tx, err := sqldb.Begin()
// 	siutils.AssertNilAndFilaNow(err, t)

// 	sqltx := siwrapper.NewSqlTx(tx)
// 	defer sqltx.Rollback()

// 	insertQuery := `
// 		insert into m_user(id, email_address, mobile_number, user_id, user_name)
// 		values((select coalesce(max(id), 0)+1 from m_user), $1, $2, $3, $4)
// 	`
// 	sqlResult, err := sqltx.Exec(insertQuery, "new@mail.com", "000", "id", "name")
// 	siutils.AssertNilAndFilaNow(err, t)

// 	sqltx.Commit()

// 	n, err := sqlResult.RowsAffected()
// 	siutils.AssertNilAndFilaNow(err, t)

// 	assert.Equal(t, int64(1), n)
// }

func TestSqlDB_QueryRow(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

	sqldb := siwrapper.NewSqlDB(db)
	// sqldb.AddSqlColumn(sicore.SqlColumn{"decimal_", sicore.SqlColTypeFloat64},
	// 	sicore.SqlColumn{"numeric_", sicore.SqlColTypeFloat64},
	// 	sicore.SqlColumn{"char_arr_", sicore.SqlColTypeUints8},
	// )

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
	siutils.NilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:12Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoAny_Slice(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

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

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:12Z"},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:13Z"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryIntoMap(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, db)

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
	siutils.NilFail(t, err)

	expected := `[{"bigint_":123,"bytea_":"MDEyMw==","char_arr_":"e2FiY2RlLGx1bmNofQ==","decimal_":123,"int2_":123,"nil":null,"numeric_":123,"str":"123","time_":"2022-01-01T12:12:12Z","varchar_arr_":"e2FiY2RlLGx1bmNofQ=="}]`
	mb, _ := json.Marshal(m)
	fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}
