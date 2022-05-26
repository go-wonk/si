package sicore

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"
)

type rowScanner struct {
	sqlColLock sync.RWMutex
	sqlCol     map[string]any
}

func newRowScanner() *rowScanner {
	return &rowScanner{
		sqlCol: make(map[string]any),
	}
}

func (rs *rowScanner) SetSqlColumn(name string, typ any) {
	rs.sqlColLock.Lock()
	defer rs.sqlColLock.Unlock()

	rs.sqlCol[name] = typ
}

func (rs *rowScanner) GetSqlColumn(name string) (any, bool) {
	rs.sqlColLock.RLock()
	defer rs.sqlColLock.RUnlock()

	if v, ok := rs.sqlCol[name]; ok {
		return v, ok
	}
	return nil, false
}

func (rs *rowScanner) ScanTypes(rows *sql.Rows, sc ...SqlColumn) ([]interface{}, []string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	for _, c := range sc {
		c.SetType(rs)
	}

	scannedRow := make([]interface{}, len(columns))

	rs.scanTypes(scannedRow, columnTypes, columns)
	return scannedRow, columns, nil
}

func (rs *rowScanner) ScanValuesIntoMap(columns []string, values []interface{}, dest *map[string]interface{}) {
	// scanIntoMap(columns, values, dest)
	rs.scanValuesIntoMap(columns, values, dest)
}

func (rs *rowScanner) scanTypes(values []interface{}, columnTypes []*sql.ColumnType, columns []string) {
	for i, ct := range columnTypes {
		if c, ok := rs.GetSqlColumn(columns[i]); ok {
			values[i] = reflect.New(reflect.PtrTo(reflect.TypeOf(c))).Interface()
		} else {
			if ct.ScanType() == nil {
				values[i] = new(interface{})
				continue
			}
			var t interface{} = reflect.New(reflect.PtrTo(ct.ScanType())).Interface()
			values[i] = t
		}
	}
}

func (rs *rowScanner) scanValuesIntoMap(columns []string, values []interface{}, dest *map[string]interface{}) {
	for idx, v := range values {
		if rv := reflect.Indirect(reflect.Indirect(reflect.ValueOf(v))); rv.IsValid() {
			(*dest)[columns[idx]] = rv.Interface()

			if valuer, ok := (*dest)[columns[idx]].(driver.Valuer); ok {
				(*dest)[columns[idx]], _ = valuer.Value()
			} else if b, ok := (*dest)[columns[idx]].(sql.RawBytes); ok {
				(*dest)[columns[idx]] = string(b)
			}
		} else {
			(*dest)[columns[idx]] = nil
		}
	}
}

func (rs *rowScanner) Scan(rows *sql.Rows, output *[]map[string]interface{}, sc ...SqlColumn) (int, error) {
	scannedRow, columns, err := rs.ScanTypes(rows, sc...)
	if err != nil {
		return 0, err
	}

	n := 0
	for rows.Next() {
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		m := make(map[string]interface{})
		rs.ScanValuesIntoMap(columns, scannedRow, &m)
		*output = append(*output, m)
		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}
