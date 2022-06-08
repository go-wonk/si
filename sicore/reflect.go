package sicore

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func findFieldNameByTag(tagKey string, t reflect.StructTag) (string, error) {
	if jt, ok := t.Lookup(tagKey); ok {
		return strings.Split(jt, ",")[0], nil
	}
	return "", fmt.Errorf("tag provided does not define a %s tag", tagKey)
}

func buildFieldNameMapByTag(structValue reflect.Value) map[string]int {
	// need to handle embedded field
	fieldNameMap := map[string]int{}
	for i := 0; i < structValue.NumField(); i++ {
		typeField := structValue.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findFieldNameByTag("json", tag)
		fieldNameMap[jname] = i
	}
	return fieldNameMap
}

func buildScanDest(columns []string, fieldNameMap map[string]int,
	structValue reflect.Value) []interface{} {

	scannedRow := make([]interface{}, len(columns))
	for i, col := range columns {
		// need to find embedded
		fieldIndex, ok := fieldNameMap[col]
		if !ok {
			// found no field corresponding to the column name
			scannedRow[i] = reflect.New(reflect.PointerTo(refTypeOfRawBytes)).Interface()
			continue
		}
		field := structValue.Field(fieldIndex)
		fieldType := field.Type()

		// this is to scan into the field directly, but it cannot handle nil
		// scannedRow[i] = field.Addr().Interface()

		switch fieldType.Kind() {
		case reflect.Pointer:
			// if a field is pointer
			scannedRow[i] = reflect.New(fieldType).Interface()
		default:
			scannedRow[i] = reflect.New(reflect.PointerTo(fieldType)).Interface()
		}
	}

	return scannedRow

}

func setScannedValue(structValue reflect.Value, scannedRow []interface{}, columns []string, fieldNameMap map[string]int) {
	// set values to the struct fields
	for i, _ := range scannedRow {
		fieldIndex := fieldNameMap[columns[i]]
		field := structValue.Field(fieldIndex)
		fieldType := field.Type()

		// skip any invalid(nil) values
		if refValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(scannedRow[i]))); refValue.IsValid() {
			switch fieldType.Kind() {
			case reflect.Pointer:
				if fieldType.Elem().Kind() == reflect.Struct {
					// embedded field
				} else {
					// field.Set(reflect.Indirect(reflect.ValueOf(scannedRow[i])))
					field.Set(refValue.Addr())
				}
			// case reflect.Int:
			// 	field.SetInt(refValue.Int())
			default:
				field.Set(refValue)
			}
		}
	}
}

// getReflectValuePointer returns a value that v points to.
// It returns error if v is not a pointer
func getReflectValuePointer(v any) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return reflect.Value{}, errors.New("not a pointer")
	}
	return reflect.Indirect(rv), nil
}

func getSliceElement(sliceValue reflect.Value) (reflect.Value, error) {
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return reflect.Value{}, errors.New("sliceValue is not slice or array")
	}

	// this only works when the slice's element is a struct
	// var elem reflect.Value
	// rvTypeElem := rv.Type().Elem()
	// switch rvTypeElem.Kind() {
	// case reflect.Pointer:
	// 	if rvTypeElem.Elem().Kind() == reflect.Struct {
	// 		elem = reflect.New(rvTypeElem.Elem())
	// 	}
	// case reflect.Struct:
	// 	elem = reflect.New(rvTypeElem).Elem()
	// }

	// to handle all kinds(struct, int, string...)
	var elem reflect.Value
	rvTypeElem := sliceValue.Type().Elem()
	switch rvTypeElem.Kind() {
	case reflect.Pointer:
		elem = reflect.New(rvTypeElem.Elem())
	default:
		elem = reflect.New(rvTypeElem).Elem()
	}

	initializeFields(elem)
	return elem, nil
}

func initializeFields(field reflect.Value) {
	n := field.NumField()
	for i := 0; i < n; i++ {
		structField := field.Type().Field(i)
		if !structField.IsExported() {
			continue
		}

		field := field.Field(i)

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

		// if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.Struct {
		// 	switch field.Interface().(type) {
		// 	// case *time.Time:
		// 	// 	// do nothing
		// 	default:
		// 		if field.IsNil() {
		// 			field.Set(reflect.New(field.Type().Elem()))
		// 		}

		// 		initializeFields(field.Elem())
		// 	}
		// } else if field.Type().Kind() == reflect.Struct {
		// 	switch field.Interface().(type) {
		// 	// case time.Time:
		// 	// 	// do nothing
		// 	default:
		// 		initializeFields(field)
		// 	}
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

func traverseFields(parent traversedField, result *[]traversedField) {
	n := parent.field.NumField()
	for i := 0; i < n; i++ {
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
			*result = append(*result, traversedField{field, append(parent.indices, i)})
			continue
		}

		var fieldValue reflect.Value
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
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
			traverseFields(traversedField{fieldValue, append(parent.indices, i)}, result)
		}

		// if field.Kind() == reflect.Pointer &&
		// 	field.Type().Elem().Kind() == reflect.Struct {

		// 	switch field.Interface().(type) {
		// 	case *time.Time, *sql.NullBool, *sql.NullByte, *sql.NullFloat64, *sql.NullInt16, *sql.NullInt32, *sql.NullInt64,
		// 		*sql.NullString, *sql.NullTime:
		// 		if field.IsNil() {
		// 			field.Set(reflect.New(field.Type().Elem()))
		// 		}
		// 		*result = append(*result, traversedField{field.Elem(), append(parent.indices, i)})
		// 	default:
		// 		if field.IsNil() {
		// 			field.Set(reflect.New(field.Type().Elem()))
		// 		}
		// 		traverseFields(traversedField{field.Elem(), append(parent.indices, i)}, result)
		// 	}
		// } else if field.Type().Kind() == reflect.Struct {
		// 	switch field.Interface().(type) {
		// 	case time.Time, sql.NullBool, sql.NullByte, sql.NullFloat64, sql.NullInt16, sql.NullInt32, sql.NullInt64,
		// 		sql.NullString, sql.NullTime:
		// 		*result = append(*result, traversedField{field, append(parent.indices, i)})
		// 	default:
		// 		traverseFields(traversedField{field, append(parent.indices, i)}, result)
		// 	}
		// } else {
		// 	// handle tag key
		// 	*result = append(*result, traversedField{field, append(parent.indices, i)})
		// }
	}
}

func buildTagNameMap(root reflect.Value, tagKey string, fields []traversedField) map[string][]int {
	m := make(map[string][]int)
	for _, v := range fields {
		field := root.Type().FieldByIndex(v.indices)
		name, err := findFieldNameByTag(tagKey, field.Tag)
		if err == nil {

			_, ok := m[name]
			if !ok {
				m[name] = v.indices
			}
		}
	}

	return m
}

func buildScanDestinations(columns []string, fieldTagMap map[string][]int,
	root reflect.Value) []interface{} {

	scannedRow := make([]interface{}, len(columns))
	for i, col := range columns {
		// need to find embedded
		fieldIndex, ok := fieldTagMap[col]
		if !ok {
			// found no field corresponding to the column name
			scannedRow[i] = reflect.New(reflect.PointerTo(refTypeOfRawBytes)).Interface()
			continue
		}
		field := root.FieldByIndex(fieldIndex)
		fieldType := field.Type()

		// this is to scan into the field directly, but it cannot handle nil
		// scannedRow[i] = field.Addr().Interface()

		switch fieldType.Kind() {
		case reflect.Pointer:
			// if a field is pointer
			scannedRow[i] = reflect.New(fieldType).Interface()
		default:
			scannedRow[i] = reflect.New(reflect.PointerTo(fieldType)).Interface()
		}
	}

	return scannedRow

}

func setScannedValues(v reflect.Value, scannedRow []interface{}, columns []string, tagNameMap map[string][]int) {
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
