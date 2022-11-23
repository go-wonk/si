package sicore

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
)

// const defaultUseSqlNullType = true
const defaultTagKey = "si"

// RowScanner scans data from sql.Rows when data type is unknown.
// By default all column type is mapped with sql.NullXXX type to be safe.
// `sqlCol` is a map to assign a data type to specific column.
type RowScanner struct {
	// sqlColLock sync.RWMutex
	sqlCol map[string]any
	tagKey string
}

func newRowScanner() *RowScanner {
	return &RowScanner{
		sqlCol: make(map[string]any),
		tagKey: defaultTagKey,
	}
}

func (rs *RowScanner) Reset(opts ...RowScannerOption) {
	for k := range rs.sqlCol {
		delete(rs.sqlCol, k)
	}
	rs.tagKey = defaultTagKey
	for _, v := range opts {
		v.apply(rs)
	}
}

func (rs *RowScanner) SetSqlColumn(name string, typ any) {
	// rs.sqlColLock.Lock()
	// defer rs.sqlColLock.Unlock()

	rs.sqlCol[name] = typ
}

func (rs *RowScanner) GetSqlColumn(name string) (any, bool) {
	// rs.sqlColLock.RLock()
	// defer rs.sqlColLock.RUnlock()

	if v, ok := rs.sqlCol[name]; ok {
		return v, ok
	}
	return nil, false
}

func (rs *RowScanner) SetTagKey(key string) {
	rs.tagKey = key
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

func (rs *RowScanner) scanTypes(values []interface{}, columnTypes []*sql.ColumnType, columns []string) {
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

		// if rs.useSqlNullType {
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
		// } else {
		// 	var t interface{} = reflect.New(ct.ScanType()).Interface()
		// 	values[i] = t
		// }
	}
}

func (rs *RowScanner) setMapValues(columns []string, values []interface{}, dest map[string]interface{}) {
	for idx := range columns {
		if rv := reflect.Indirect(reflect.ValueOf(values[idx])); rv.IsValid() {
			var rvi interface{} = rv.Interface()

			switch v := rvi.(type) {
			case driver.Valuer:
				dest[columns[idx]], _ = v.Value()
			// case sql.RawBytes:
			// 	dest[columns[idx]] = string(v)
			default:
				dest[columns[idx]] = rvi
			}
		} else {
			dest[columns[idx]] = nil
		}
	}
}

// ScanMapSlice scans `rows` into `output`.
func (rs *RowScanner) ScanMapSlice(rows *sql.Rows, output *[]map[string]interface{}) (int, error) {
	scannedRow, columns, err := rs.ScanTypes(rows)
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
		rs.setMapValues(columns, scannedRow, m)

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

// ScanStructs scans `rows` into `output`. `output` should be a slice of structs.
func (rs *RowScanner) ScanStructs(rows *sql.Rows, output any) (int, error) {
	sliceValue, err := valueOfAnyPtr(output)
	if err != nil {
		return 0, err
	}

	if !isSliceKind(sliceValue) {
		return 0, errors.New("ouput is not a slice")
	}

	elemType, isPtr := typeOfSliceElem(sliceValue)

	var elemValue reflect.Value
	if isPtr {
		elemValue = newValueOfSliceElemPtr(elemType)
	} else {
		elemValue = newValueOfSliceElem(elemType)
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	for i := range columns {
		columns[i] = strings.ToLower(columns[i])
	}

	n := 0 // num rows

	var traversedFields []traversedField
	var fieldsToInitialize [][]int
	traverseFields(traversedField{elemValue, []int{}}, rs.tagKey, &traversedFields, &fieldsToInitialize)
	tagNameMap := makeNameMap(elemValue, rs.tagKey, traversedFields)

	scannedRow, err := buildDestinations(columns, tagNameMap, elemValue)
	if err != nil {
		return 0, err
	}
	for rows.Next() {

		// scan the values
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		if isPtr {
			elemValue = newValueOfSliceElemPtr(elemType)
		} else {
			elemValue = newValueOfSliceElem(elemType)
		}

		initializeFieldsWithIndices(elemValue, fieldsToInitialize)

		// set values to the struct fields
		setStructValues(elemValue, scannedRow, columns, tagNameMap)

		// append element to slice
		if isPtr {
			elemValue = elemValue.Addr()
		}
		sliceValue.Set(reflect.Append(sliceValue, elemValue))

		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// ScanStruct scans `rows` into `output`. `output` is a struct.
func (rs *RowScanner) ScanStruct(rows *sql.Rows, output any) error {
	rv, err := valueOfAnyPtr(output)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	for i := range columns {
		columns[i] = strings.ToLower(columns[i])
	}

	var traversedFields []traversedField
	var fieldsToInitialize [][]int
	traverseFields(traversedField{rv, []int{}}, rs.tagKey, &traversedFields, &fieldsToInitialize)
	tagNameMap := makeNameMap(rv, rs.tagKey, traversedFields)

	dest, err := buildDestinations(columns, tagNameMap, rv)
	if err != nil {
		return err
	}
	if rows.Next() {
		// scan the values
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}

		// set values to the struct fields
		setStructValues(rv, dest, columns, tagNameMap)
	} else {
		return sql.ErrNoRows
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// ScanStructs scans `rows` into `output`. `output` should be a slice of structs.
func (rs *RowScanner) ScanPrimary(row *sql.Row, output any) error {
	rv, err := valueOfAnyPtr(output)
	if err != nil {
		return err
	}

	dest := reflect.New(reflect.PointerTo(rv.Type())).Interface()

	// scan the values
	err = row.Scan(dest)
	if err != nil {
		return err
	}

	// set values to the struct fields
	if refValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(dest))); refValue.IsValid() {
		rv.Set(refValue)
	}

	return nil
}
