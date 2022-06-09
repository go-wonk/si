package sicore

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// findTagName finds tag name with `tagKey` in `t`
func findTagName(tagKey string, t reflect.StructTag) (string, error) {
	if jt, ok := t.Lookup(tagKey); ok {
		return strings.Split(jt, ",")[0], nil
	}
	return "", fmt.Errorf("tagKey '%s' was not found", tagKey)
}

// valueOfAnyPtr returns a value that v points to.
// It returns error if v is not a pointer.
func valueOfAnyPtr(v any) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return reflect.Value{}, errors.New("not a pointer")
	}
	return rv.Elem(), nil
}

// isSliceKind returns true if `v` is a slice or array
func isSliceKind(v reflect.Value) bool {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}
	return true
}

// typeOfSliceElem returns `elemType` that is the element type of a slice, `sv`,
// and `isPtr` is true only if `elemType` is a pointer.
func typeOfSliceElem(sv reflect.Value) (elemType reflect.Type, isPtr bool) {
	elemType = sv.Type().Elem()
	switch elemType.Kind() {
	case reflect.Pointer:
		isPtr = true
	default:
		isPtr = false
	}

	return
}

// newValueOfSliceElem creates a new element of `elemType`.
func newValueOfSliceElem(elemType reflect.Type) (elem reflect.Value) {
	elem = reflect.New(elemType).Elem()
	return
}

// newValueOfSliceElemPtr creates a new element of `elemType` when it is a pointer.
func newValueOfSliceElemPtr(elemType reflect.Type) (elem reflect.Value) {
	elem = reflect.New(elemType.Elem()).Elem()
	return
}

// initializeFields traverses all fields recursivley and initialize nil pointer struct
func initializeFields(v reflect.Value) {
	n := v.NumField()
	for i := 0; i < n; i++ {
		structField := v.Type().Field(i)
		if !structField.IsExported() {
			continue
		}

		field := v.Field(i)

		// check if a field is pointer
		var fieldTypeKind reflect.Kind
		if field.Kind() == reflect.Pointer {
			fieldTypeKind = field.Type().Elem().Kind()
		} else {
			fieldTypeKind = field.Type().Kind()
		}

		// skip if it is not a struct
		if fieldTypeKind != reflect.Struct {
			continue
		}

		// initialize a field if it is nil
		var fieldValue reflect.Value
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			fieldValue = field.Elem()
		} else {
			fieldValue = field
		}

		// initialize struct's nested field
		initializeFields(fieldValue)
	}
}

// initializeFieldsWithIndices initializes directly using `indices`
func initializeFieldsWithIndices(v reflect.Value, indices [][]int) {
	for _, s := range indices {
		field := v.FieldByIndex(s)
		// if field.Kind() == reflect.Pointer && field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
		// }
	}
}

type traversedField struct {
	field   reflect.Value
	indices []int
}

type ScanValuer interface {
	Scan(value any) error
	Value() (driver.Value, error)
}

func traverseFields(parent traversedField, result *[]traversedField, resultInitialize *[][]int) {
	n := parent.field.NumField()
	for i := 0; i < n; i++ {
		// skip any unexported(private) fields
		structField := parent.field.Type().Field(i)
		if !structField.IsExported() {
			continue
		}

		field := parent.field.Field(i)

		var fieldTypeKind reflect.Kind
		if field.Kind() == reflect.Pointer {
			fieldTypeKind = field.Type().Elem().Kind()
		} else {
			fieldTypeKind = field.Type().Kind()
		}

		if fieldTypeKind != reflect.Struct {
			if fieldTypeKind == reflect.Interface && field.NumMethod() > 0 {
				continue
			}
			*result = append(*result, traversedField{field, append(parent.indices, i)})
			continue
		}

		var fieldValue reflect.Value
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
				*resultInitialize = append(*resultInitialize, append(parent.indices, i))
			}
			fieldValue = field.Elem()
		} else {
			fieldValue = field
		}

		switch field.Interface().(type) {
		case *time.Time, *sql.NullBool, *sql.NullByte, *sql.NullFloat64, *sql.NullInt16, *sql.NullInt32, *sql.NullInt64,
			*sql.NullString, *sql.NullTime,
			time.Time, sql.NullBool, sql.NullByte, sql.NullFloat64, sql.NullInt16, sql.NullInt32, sql.NullInt64,
			sql.NullString, sql.NullTime:

			*result = append(*result, traversedField{fieldValue, append(parent.indices, i)})
		default:
			traverseFields(traversedField{fieldValue, append(parent.indices, i)}, result, resultInitialize)
		}
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func makeNameMap(root reflect.Value, tagKey string, fields []traversedField) map[string][]int {
	m := make(map[string][]int)
	for _, v := range fields {
		field := root.Type().FieldByIndex(v.indices)
		name, err := findTagName(tagKey, field.Tag)
		if err != nil {
			if len(field.Name) == 0 {
				continue
			}
			name = ToSnake(field.Name)
		}
		if len(name) == 0 {
			continue
		}
		_, ok := m[name]
		if !ok {
			m[name] = v.indices
		}
	}

	return m
}

func buildDestinations(columns []string, fieldTagMap map[string][]int, root reflect.Value) ([]interface{}, error) {

	dest := make([]interface{}, len(columns))
	for i, col := range columns {
		// need to find embedded
		fieldIndex, ok := fieldTagMap[col]
		if !ok {
			// found no field corresponding to the column name

			// proceed even if selected columns are not matched with struct
			// scannedRow[i] = reflect.New(reflect.PointerTo(refTypeOfRawBytes)).Interface()
			// continue

			return nil, fmt.Errorf("column '%s' was not found", col)
		}
		field := root.FieldByIndex(fieldIndex)
		fieldType := field.Type()

		// this is to scan into the field directly, but it cannot handle nil
		// scannedRow[i] = field.Addr().Interface()

		switch fieldType.Kind() {
		case reflect.Pointer:
			// if a field is pointer
			dest[i] = reflect.New(fieldType).Interface()
		default:
			dest[i] = reflect.New(reflect.PointerTo(fieldType)).Interface()
		}
	}

	return dest, nil

}

func setStructValues(v reflect.Value, scannedRow []interface{}, columns []string, tagNameMap map[string][]int) {
	// set values to the struct fields
	for i := range scannedRow {
		indices, ok := tagNameMap[columns[i]]
		if !ok {
			continue
		}
		field := v.FieldByIndex(indices)
		fieldType := field.Type()

		// skip any invalid(nil) values, so skipped fields will have their default values like 0, "", false and etc.
		if refValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(scannedRow[i]))); refValue.IsValid() {
			switch fieldType.Kind() {
			case reflect.Pointer:
				field.Set(refValue.Addr())
			default:
				field.Set(refValue)
			}
		}
	}
}
