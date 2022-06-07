package sicore

import (
	"database/sql"
	"database/sql/driver"
	"errors"
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

func (rs *rowScanner) scanTypes(values []interface{}, columnTypes []*sql.ColumnType, columns []string) {
	for i, ct := range columnTypes {
		if len(rs.sqlCol) > 0 {
			if c, ok := rs.GetSqlColumn(columns[i]); ok {
				values[i] = reflect.New(reflect.TypeOf(c)).Interface()
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
				values[i] = reflect.New(refTypeOfRawBytes).Interface()
			case refTypeOfBytesTypeValue:
				values[i] = reflect.New(refTypeOfRawBytes).Interface()
			case refTypeOfByteTypeValue:
				values[i] = reflect.New(refTypeOfNullByte).Interface()
			case refTypeOfBoolTypeValue:
				values[i] = reflect.New(refTypeOfNullBool).Interface()
			case refTypeOfStringTypeValue:
				values[i] = reflect.New(refTypeOfNullString).Interface()
			case refTypeOfFloat32TypeValue:
				values[i] = reflect.New(refTypeOfNullFloat32).Interface()
			case refTypeOfFloat64TypeValue:
				values[i] = reflect.New(refTypeOfNullFloat64).Interface()
			case refTypeOfIntTypeValue:
				values[i] = reflect.New(refTypeOfNullInt).Interface()
			case refTypeOfInt8TypeValue:
				values[i] = reflect.New(refTypeOfNullInt8).Interface()
			case refTypeOfInt16TypeValue:
				values[i] = reflect.New(refTypeOfNullInt16).Interface()
			case refTypeOfInt32TypeValue:
				values[i] = reflect.New(refTypeOfNullInt32).Interface()
			case refTypeOfInt64TypeValue:
				values[i] = reflect.New(refTypeOfNullInt64).Interface()
			case refTypeOfUintTypeValue:
				values[i] = reflect.New(refTypeOfNullUint).Interface()
			case refTypeOfUint8TypeValue:
				values[i] = reflect.New(refTypeOfNullUint8).Interface()
			case refTypeOfUint16TypeValue:
				values[i] = reflect.New(refTypeOfNullUint16).Interface()
			case refTypeOfUint32TypeValue:
				values[i] = reflect.New(refTypeOfNullUint32).Interface()
			case refTypeOfUint64TypeValue:
				values[i] = reflect.New(refTypeOfNullUint64).Interface()
			case refTypeOfTimeTypeValue:
				values[i] = reflect.New(refTypeOfNullTime).Interface()
			default:
				switch ct.DatabaseTypeName() {
				case "NUMERIC", "DECIMAL", "NUMBER":
					values[i] = reflect.New(refTypeOfNullFloat64).Interface()
				case "VARCHAR", "VARCHAR2", "NVARCHAR", "CHAR", "NCHAR", "TEXT":
					values[i] = reflect.New(refTypeOfNullString).Interface()
				default:
					var t interface{} = reflect.New(ct.ScanType()).Interface()
					values[i] = t
				}
			}
		} else {
			var t interface{} = reflect.New(ct.ScanType()).Interface()
			values[i] = t
		}
	}
}

func (rs *rowScanner) scanValuesMap(columns []string, values []interface{}, dest map[string]interface{}) {
	for idx := range columns {
		if rv := reflect.Indirect(reflect.ValueOf(values[idx])); rv.IsValid() {
			var rvi interface{} = rv.Interface()

			switch v := rvi.(type) {
			case driver.Valuer:
				dest[columns[idx]], _ = v.Value()
			case sql.RawBytes:
				dest[columns[idx]] = string(v)
			default:
				dest[columns[idx]] = rvi
			}
		} else {
			dest[columns[idx]] = nil
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
		rs.scanValuesMap(columns, scannedRow, m)

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
func (rs *rowScanner) ScanStructs_Deprecated(rows *sql.Rows, output any, sc ...SqlColumn) (int, error) {
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

// Scan scans rows' data type into a slice of interface{} first, then read actual values from rows into the slice
func (rs *rowScanner) ScanStructs2_Deprecated(rows *sql.Rows, output any, sc ...SqlColumn) (int, error) {
	rv := reflect.Indirect(reflect.ValueOf(output))
	if rv.Kind() != reflect.Slice {
		return 0, errors.New("not slice")
	}
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	n := 0

	foundFieldNames := false
	var fieldNameMap map[string]int

	scanned := false
	var scannedRow []interface{}

	rvTypeElem := rv.Type().Elem() // type of the slice element
	rvTypeElemKind := rvTypeElem.Kind()
	for rows.Next() {
		var elem reflect.Value // slice's element
		if rvTypeElemKind == reflect.Pointer {
			elem = reflect.New(rvTypeElem.Elem())
		} else if rvTypeElemKind == reflect.Struct {
			elem = reflect.New(rvTypeElem).Elem()
		}

		if !foundFieldNames {
			fieldNameMap = buildFieldNameMapByTag(elem)
			foundFieldNames = true
		}

		if !scanned {
			scannedRow = buildScanDest(columns, fieldNameMap, elem)
			scanned = true
		}

		// scan the values
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		// set values to the struct fields
		setScannedValue(elem, scannedRow, columns, fieldNameMap)

		// append element to slice
		rv.Set(reflect.Append(rv, elem))

		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Scan scans rows' data type into a slice of interface{} first, then read actual values from rows into the slice
func (rs *rowScanner) ScanStructs(rows *sql.Rows, output any, sc ...SqlColumn) (int, error) {
	sliceValue, err := getReflectValuePointer(output)
	if err != nil {
		return 0, err
	}

	elemValue, err := getSliceElement(sliceValue)
	if err != nil {
		return 0, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	n := 0 // num rows

	var traversedFields []traversedField
	traverseFields(traversedField{elemValue, []int{}}, &traversedFields)
	tagNameMap := buildTagNameMap(elemValue, "json", traversedFields)

	scannedRow := buildScanDestinations(columns, tagNameMap, elemValue)
	for rows.Next() {

		// scan the values
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		elem, err := getSliceElement(sliceValue)
		if err != nil {
			return 0, err
		}
		// set values to the struct fields
		setScannedValues(elem, scannedRow, columns, tagNameMap)

		// append element to slice
		sliceValue.Set(reflect.Append(sliceValue, elem))

		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}
