package sicore

import (
	"fmt"
	"reflect"
	"strings"
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

func instantiateStructField(field reflect.Value) {
	if field.Kind() == reflect.Pointer &&
		field.Type().Elem().Kind() == reflect.Struct &&
		field.IsNil() {

		field.Set(reflect.New(field.Type().Elem()))
	}
}

type traversedField struct {
	field   reflect.Value
	indices []int
}

func traverseFields(tagKey string, parent traversedField, result *[]traversedField) {

	n := parent.field.NumField()
	// fmt.Println(n)
	for i := 0; i < n; i++ {
		field := parent.field.Field(i)
		if field.Kind() == reflect.Pointer &&
			field.Type().Elem().Kind() == reflect.Struct {

			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}

			traverseFields(tagKey, traversedField{field.Elem(), append(parent.indices, i)}, result)
		} else if field.Type().Kind() == reflect.Struct {
			traverseFields(tagKey, traversedField{field, append(parent.indices, i)}, result)
		} else {
			// handle tag key
			*result = append(*result, traversedField{field, append(parent.indices, i)})
			// fmt.Println(parent.field.Field(i).Type(), toTraverse{field, append(parent.indices, i)})
		}
	}
}
