package sisql_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-wonk/si/sisql"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
	"github.com/stretchr/testify/assert"
)

func TestSqlDB_QueryRow(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db)

	query := `
		select 
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
	`

	var tim sql.NullTime
	row := sqldb.QueryRow(query)
	err := row.Scan(&tim)
	siutils.AssertNilFail(t, err)

	expected := `2022-01-01 12:12:12 +0000 UTC`
	assert.Equal(t, expected, tim.Time.String())
}

func TestSqlDB_QueryRowStruct(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select id, name, email_address, borrowed, 23 as book_id from student order by id limit 1
	`

	// tl := Table{}
	var tl testmodels.Student
	err := sqldb.QueryRowStruct(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `{"id":1,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23}`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDB_QueryRowPrimary(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select 12 as id
	`

	var tl int
	err := sqldb.QueryRowPrimary(query, &tl)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, 12, tl)
}

func TestSqlDBQueryStructsBasicDataType(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select 0 as bool_value, '' as string_value, 
			0 as int_value, 8 as int8_value, 16 as int16_value, 32 as int32_value, 64 as int64_value,
			0 as uint_value, 8 as uint8_value, 16 as uint16_value, 32 as uint32_value, 64 as uint64_value,
			41 as byte_value, 
			32.32 as float32_value, 64.64 as float64_value,
			'bytes1234'::bytea as bytes_value,
			1234.1234::decimal(20,4) as big_float_value,
			0 as bool_ptr_value, '' as string_ptr_value, 
			0 as int_ptr_value, 8 as int8_ptr_value, 16 as int16_ptr_value, 32 as int32_ptr_value, 64 as int64_ptr_value,
			0 as uint_ptr_value, 8 as uint8_ptr_value, 16 as uint16_ptr_value, 32 as uint32_ptr_value, 64 as uint64_ptr_value,
			41 as byte_ptr_value, 
			32.32 as float32_ptr_value, 64.64 as float64_ptr_value,
			'bytes1234'::bytea as bytes_ptr_value,
			1234.1234::decimal(20,4) as big_float_ptr_value
	`

	var l testmodels.BasicDataTypeList
	_, err := sqldb.QueryStructs(query, &l)
	siutils.AssertNilFail(t, err)

	// jsonValue := l.String()
	// fmt.Println(jsonValue)

	expected := `[{"bool_value":false,"string_value":"","int_value":0,"int8_value":8,"int16_value":16,"int32_value":32,"int64_value":64,"uint_value":0,"uint8_value":8,"uint16_value":16,"uint32_value":32,"uint64_value":64,"byte_value":41,"float32_value":32.32,"float64_value":64.64,"bytes_value":"Ynl0ZXMxMjM0","big_float_value":"MTIzNC4xMjM0","bool_ptr_value":false,"string_ptr_value":"","int_ptr_value":0,"int8_ptr_value":8,"int16_ptr_value":16,"int32_ptr_value":32,"int64_ptr_value":64,"uint_ptr_value":0,"uint8_ptr_value":8,"uint16_ptr_value":16,"uint32_ptr_value":32,"uint64_ptr_value":64,"byte_ptr_value":41,"float32_ptr_value":32.32,"float64_ptr_value":64.64,"bytes_ptr_value":"Ynl0ZXMxMjM0","big_float_ptr_value":"MTIzNC4xMjM0"}]`
	assert.Equal(t, expected, l.String())

}

func TestSqlDBQueryStructsBasicDataTypeLimit(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	// float values can only express 17 digits from the left, the rest is zero'ed
	f64 := 1.797693134862315708145274237317043567981e+308
	fmt.Println(f64) // 1.7976931348623157e+308

	f64 = 12345678901234567890.123456789
	fmt.Println(f64) // 1.2345678901234567e+19

	query := `
		select 0 as bool_value, 'this is a string value' as string_value, 
			9223372036854775807 as int_value, 127 as int8_value, 32767 as int16_value, 2147483647 as int32_value, 9223372036854775807 as int64_value,
			18446744073709551615 as uint_value, 255 as uint8_value, 65535 as uint16_value, 4294967295 as uint32_value, 18446744073709551615 as uint64_value,
			41 as byte_value, 
			32.32 as float32_value, 12345678901234567890.1234::varchar(50) as float64_value
	`

	var l testmodels.BasicDataTypeList
	_, err := sqldb.QueryStructs(query, &l)
	siutils.AssertNilFail(t, err)

	// jsonValue := l.String()
	// fmt.Println(l.String())

	expected := `[{"bool_value":false,"string_value":"this is a string value","int_value":9223372036854775807,"int8_value":127,"int16_value":32767,"int32_value":2147483647,"int64_value":9223372036854775807,"uint_value":18446744073709551615,"uint8_value":255,"uint16_value":65535,"uint32_value":4294967295,"uint64_value":18446744073709551615,"byte_value":41,"float32_value":32.32,"float64_value":12345678901234567000,"bytes_value":null,"big_float_value":null,"bool_ptr_value":null,"string_ptr_value":null,"int_ptr_value":null,"int8_ptr_value":null,"int16_ptr_value":null,"int32_ptr_value":null,"int64_ptr_value":null,"uint_ptr_value":null,"uint8_ptr_value":null,"uint16_ptr_value":null,"uint32_ptr_value":null,"uint64_ptr_value":null,"byte_ptr_value":null,"float32_ptr_value":null,"float64_ptr_value":null,"bytes_ptr_value":null,"big_float_ptr_value":null}]`
	assert.Equal(t, expected, l.String())

}

func TestSqlDBQueryStructsBasicDataTypeBig(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select '12345678901234567890.1234'::bytea as bytes_value,
			'12345678901234567890.1234'::bytea as bytes_ptr_value,
			12345678901234567890.1234::decimal(24,4) as big_float_value,
			12345678901234567890.1234::decimal(24,4) as big_float_ptr_value
	`

	var l testmodels.BasicDataTypeList
	_, err := sqldb.QueryStructs(query, &l)
	siutils.AssertNilFail(t, err)

	// jsonValue := l.String()
	// fmt.Println(l.String())

	expected := `[{"bool_value":false,"string_value":"","int_value":0,"int8_value":0,"int16_value":0,"int32_value":0,"int64_value":0,"uint_value":0,"uint8_value":0,"uint16_value":0,"uint32_value":0,"uint64_value":0,"byte_value":0,"float32_value":0,"float64_value":0,"bytes_value":"MTIzNDU2Nzg5MDEyMzQ1Njc4OTAuMTIzNA==","big_float_value":"MTIzNDU2Nzg5MDEyMzQ1Njc4OTAuMTIzNA==","bool_ptr_value":null,"string_ptr_value":null,"int_ptr_value":null,"int8_ptr_value":null,"int16_ptr_value":null,"int32_ptr_value":null,"int64_ptr_value":null,"uint_ptr_value":null,"uint8_ptr_value":null,"uint16_ptr_value":null,"uint32_ptr_value":null,"uint64_ptr_value":null,"byte_ptr_value":null,"float32_ptr_value":null,"float64_ptr_value":null,"bytes_ptr_value":"MTIzNDU2Nzg5MDEyMzQ1Njc4OTAuMTIzNA==","big_float_ptr_value":"MTIzNDU2Nzg5MDEyMzQ1Njc4OTAuMTIzNA=="}]`
	assert.Equal(t, expected, l.String())

}

func TestSqlDBQueryStructsSimple(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select id, name, email_address, borrowed, 23 as book_id from student order by id limit 10
	`

	// tl := Table{}
	var tl testmodels.StudentList
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	// expected := `[{"id":1,"email_address":"asdf","name":"asdf","borrowed":false,"book_id":23},{"id":2,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":3,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":4,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":5,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":6,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":7,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":8,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":9,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23},{"id":10,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23}]`
	// assert.Equal(t, expected, tl.String())
}

func TestSqlDBQueryStructsNil(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	// te := &testmodels.Embedded{
	// 	Nil: "not nil embedded",
	// }
	// ts := &testmodels.Sample{
	// 	te, "", 1, nil,
	// }
	// fmt.Println(ts.String())

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select null as nil, '2' as embedded_nil, 123::integer as int2_, 234::integer as int3_
		union all
		select 'not null' as nil, '3' as embedded_nil, null as int2_, null as int3_
		union all
		select 'not null' as nil, null as embedded_nil, null as int2_, null as int3_
		union all
		select 'not null' as nil, null as embedded_nil, 99999::integer as int2_, 88888::integer as int3_
	`

	// tl := Table{}
	var tl testmodels.SampleList
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)
	// fmt.Println(tl.String())

	expected := `[{"embedded_nil_":"","embedded_nil":"2","nil":"","int2_":123,"int3_":234},{"embedded_nil_":"","embedded_nil":"3","nil":"not null","int2_":0,"int3_":null},{"embedded_nil_":"","embedded_nil":"","nil":"not null","int2_":0,"int3_":null},{"embedded_nil_":"","embedded_nil":"","nil":"not null","int2_":99999,"int3_":88888}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDBQueryStructs(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select null as nil, 
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_ptr_,
			'child of somebody' as child_name,
			't string' as table_str,
			't string 2' as table_str2,
			0 as table_interface,
			0 as any_value,
			'asdf' as null_string,
			'asdf-ptr' as null_string_ptr
	`

	tl := testmodels.TableList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:12Z","time_ptr_":"2022-01-01T12:12:12Z","c3":{"child_name":"child of somebody"},"c4":{"child_name":""},"table_str":"t string","table_str2":"t string 2","table_interface":0,"Tabler":null,"any_value":0,"null_string":{"String":"asdf","Valid":true},"null_string_ptr":{"String":"asdf-ptr","Valid":true}}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDBQueryStructs2Rows(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select null as nil, 
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
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_, 
			123456789123::bigint as bigint_, 
			'{"abcde", "lunch"}'::char(5)[] as char_arr_, 
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121213', 'YYYYMMDDHH24MISS') as time_
	`

	tl := testmodels.TableList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:12Z","time_ptr_":"0001-01-01T00:00:00Z","c3":{"child_name":""},"c4":{"child_name":""},"table_str":"","table_str2":null,"table_interface":null,"Tabler":null,"any_value":null,"null_string":{"String":"","Valid":false},"null_string_ptr":{"String":"","Valid":false}},{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123456789123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"MDEyMw==","time_":"2022-01-01T12:12:13Z","time_ptr_":"0001-01-01T00:00:00Z","c3":{"child_name":""},"c4":{"child_name":""},"table_str":"","table_str2":null,"table_interface":null,"Tabler":null,"any_value":null,"null_string":{"String":"","Valid":false},"null_string_ptr":{"String":"","Valid":false}}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDBQueryMaps(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db)

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

	expected := `[{"bigint_":123,"bytea_":"MDEyMw==","char_arr_":"e2FiY2RlLGx1bmNofQ==","decimal_":123,"int2_":123,"nil":null,"numeric_":123,"str":"123","time_":"2022-01-01T12:12:12Z","varchar_arr_":"e2FiY2RlLGx1bmNofQ=="}]`
	mb, _ := json.Marshal(m)
	// fmt.Println(string(mb))
	assert.Equal(t, expected, string(mb))

}

func TestSqlDBQueryMapsBool(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

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

func TestSqlDBQueryMapsBoolWithSqlColumn(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTypedBool("true_1").WithTypedBool("true_2").WithTypedBool("false_1").WithTypedBool("false_2")
	// sicore.SqlColumn{Name: "true_1", Type: sicore.SqlColTypeBool},
	// sicore.SqlColumn{Name: "true_2", Type: sicore.SqlColTypeBool},
	// sicore.SqlColumn{Name: "false_1", Type: sicore.SqlColTypeBool},
	// sicore.SqlColumn{Name: "false_2", Type: sicore.SqlColTypeBool},

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

func TestSqlDBQueryStructsNoTag(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select null as nil_value, 
			123::integer as int_value,
			123::decimal(24,4) as decimal_value,
			'some string' as some_string_value
	`

	tl := testmodels.TableWithNoTagList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil_value":"","IntValue":123,"DecimalValue":123,"SomeStringValue":"some string"}]`
	assert.Equal(t, expected, tl.String())
}

func TestSqlDBQueryStructsPtrElem(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, db)

	sqldb := sisql.NewSqlDB(db).WithTagKey("json")

	query := `
		select null as nil_value, 
			123::integer as int_value,
			123::decimal(24,4) as decimal_value,
			'some string' as some_string_value
		union all
		select 'not null' as nil_value, 
			987::integer as int_value,
			654::decimal(24,4) as decimal_value,
			'2some string2' as some_string_value
	`

	tl := testmodels.TableWithPtrElemList{}
	_, err := sqldb.QueryStructs(query, &tl)
	siutils.AssertNilFail(t, err)

	expected := `[{"nil_value":"","IntValue":123,"DecimalValue":123,"SomeStringValue":"some string"},{"nil_value":"not null","IntValue":987,"DecimalValue":654,"SomeStringValue":"2some string2"}]`
	assert.Equal(t, expected, tl.String())
}
