package sicore

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"
)

type RowScanner struct {
	l      sync.RWMutex
	sqlCol map[string]any
}

func NewRowScanner() *RowScanner {
	return &RowScanner{
		sqlCol: make(map[string]any),
	}
}

func (rs *RowScanner) SetSqlColumn(name string, typ any) {
	rs.l.Lock()
	defer rs.l.Unlock()

	rs.sqlCol[name] = typ
}

func (rs *RowScanner) GetSqlColumn(name string) (any, bool) {
	rs.l.RLock()
	defer rs.l.RUnlock()

	if v, ok := rs.sqlCol[name]; ok {
		return v, ok
	}
	return nil, false
}

func (rs *RowScanner) ScanTypes(rows *sql.Rows) ([]interface{}, []string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	scannedRow := make([]interface{}, len(columns))

	rs.scanTypes(scannedRow, columnTypes, columns)
	return scannedRow, columns, nil
}

func (rs *RowScanner) ScanValuesIntoMap(columns []string, values []interface{}, dest *map[string]interface{}) {
	// scanIntoMap(columns, values, dest)
	rs.scanValuesIntoMap(columns, values, dest)
}

func (rs *RowScanner) scanTypes(values []interface{}, columnTypes []*sql.ColumnType, columns []string) {
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

func (rs *RowScanner) scanValuesIntoMap(columns []string, values []interface{}, dest *map[string]interface{}) {
	for idx, v := range values {
		if rv := reflect.Indirect(reflect.Indirect(reflect.ValueOf(v))); rv.IsValid() {
			(*dest)[columns[idx]] = rv.Interface()
			// fmt.Printf("%v: %p, %T, %v\n", columns[idx], (*dest)[columns[idx]], (*dest)[columns[idx]], (*dest)[columns[idx]])

			if valuer, ok := (*dest)[columns[idx]].(driver.Valuer); ok {
				(*dest)[columns[idx]], _ = valuer.Value()
			} else if b, ok := (*dest)[columns[idx]].(sql.RawBytes); ok {
				(*dest)[columns[idx]] = string(b)
			} else {
				(*dest)[columns[idx]] = rv.Interface()
			}
		} else {
			(*dest)[columns[idx]] = nil
		}
	}
}
