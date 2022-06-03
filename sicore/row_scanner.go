package sicore

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"

	"github.com/go-wonk/si/siutils"
)

const defaultUseSqlNullType = true

// rowScanner scans data from sql.Rows when data type is unknown.
// By default all column type is mapped with sql.NullXXX type to be safe.
// `sqlCol` is a map to assign a data type to specific column.
type rowScanner struct {
	sqlColLock     sync.RWMutex
	sqlCol         map[string]any
	useSqlNullType bool
}

func newRowScanner() *rowScanner {
	return &rowScanner{
		sqlCol:         make(map[string]any),
		useSqlNullType: defaultUseSqlNullType,
	}
}

func (rs *rowScanner) Reset(useSqlNullType bool) {
	for k := range rs.sqlCol {
		delete(rs.sqlCol, k)
	}
	rs.useSqlNullType = useSqlNullType
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
		if len(rs.sqlCol) > 0 {
			if c, ok := rs.GetSqlColumn(columns[i]); ok {
				values[i] = reflect.New(reflect.PtrTo(reflect.TypeOf(c))).Interface()
				continue
			}
		}

		if ct.ScanType() == nil {
			values[i] = new(interface{})
			continue
		}

		if rs.useSqlNullType {
			switch ct.ScanType() {
			case refTypeOfRawBytes:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfRawBytes)).Interface()
			case refTypeOfBytesTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfRawBytes)).Interface()
			case refTypeOfByteTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullByte)).Interface()
			case refTypeOfBoolTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullBool)).Interface()
			case refTypeOfStringTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullString)).Interface()
			case refTypeOfFloat32TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullFloat32)).Interface()
			case refTypeOfFloat64TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullFloat64)).Interface()
			case refTypeOfIntTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullInt)).Interface()
			case refTypeOfInt8TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullInt8)).Interface()
			case refTypeOfInt16TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullInt16)).Interface()
			case refTypeOfInt32TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullInt32)).Interface()
			case refTypeOfInt64TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullInt64)).Interface()
			case refTypeOfUintTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullUint)).Interface()
			case refTypeOfUint8TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullUint8)).Interface()
			case refTypeOfUint16TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullUint16)).Interface()
			case refTypeOfUint32TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullUint32)).Interface()
			case refTypeOfUint64TypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullUint64)).Interface()
			case refTypeOfTimeTypeValue:
				values[i] = reflect.New(reflect.PtrTo(refTypeOfNullTime)).Interface()
			default:
				switch ct.DatabaseTypeName() {
				case "NUMERIC", "DECIMAL", "NUMBER":
					values[i] = reflect.New(reflect.PtrTo(refTypeOfNullFloat64)).Interface()
				case "VARCHAR", "VARCHAR2", "NVARCHAR", "CHAR", "NCHAR", "TEXT":
					values[i] = reflect.New(reflect.PtrTo(refTypeOfNullString)).Interface()
				default:
					var t interface{} = reflect.New(reflect.PtrTo(ct.ScanType())).Interface()
					values[i] = t
				}
			}
		} else {
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

func (rs *rowScanner) scanValuesMap(columns []string, values []interface{}, dest map[string]interface{}) {
	for idx, v := range values {
		columnName := columns[idx]
		if rv := reflect.Indirect(reflect.Indirect(reflect.ValueOf(v))); rv.IsValid() {
			var rvi interface{} = rv.Interface()

			if valuer, ok := rvi.(driver.Valuer); ok {
				dest[columnName], _ = valuer.Value()
			} else if b, ok := rvi.(sql.RawBytes); ok {
				dest[columnName] = string(b)
			} else {
				dest[columnName] = rvi
			}
		} else {
			dest[columnName] = nil
		}
	}
}

// Scan scans rows' data type into a slice of interface{} first, then read actual values from rows into the slice
func (rs *rowScanner) Scan(rows *sql.Rows, output *[]map[string]interface{}, sc ...SqlColumn) (int, error) {
	scannedRow, columns, err := rs.ScanTypes(rows, sc...)
	if err != nil {
		return 0, err
	}

	n := 0
	numCol := len(columns)
	for rows.Next() {
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		m := make(map[string]interface{}, numCol)
		rs.ScanValuesIntoMap(columns, scannedRow, &m)

		*output = append(*output, m)
		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	*output = (*output)[:n]
	return n, nil
}

// Scan scans rows' data type into a slice of interface{} first, then read actual values from rows into the slice
func (rs *rowScanner) ScanStructs(rows *sql.Rows, output any, sc ...SqlColumn) (int, error) {
	scannedRow, columns, err := rs.ScanTypes(rows, sc...)
	if err != nil {
		return 0, err
	}

	ms := getMapSlice()
	defer putMapSlice(ms)

	n := 0
	_ = len(columns)
	for rows.Next() {
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		if len(ms) < n+1 {
			_, err := growMapSlice(&ms, 100)
			if err != nil {
				return 0, err
			}
		}
		makeMapIfNil(&ms[n])
		rs.scanValuesMap(columns, scannedRow, ms[n])

		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	// simple, not very ideal json unmarshal
	if err = siutils.DecodeAny(ms[:n], output); err != nil {
		return 0, err
	}
	return n, nil
}
