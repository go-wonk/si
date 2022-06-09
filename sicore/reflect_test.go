package sicore

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-wonk/si/siutils"
)

type TestEmbeddedLevel2 struct {
	Name string `json:"name_level2"`
}
type TestEmbeddedLevel1 struct {
	*TestEmbeddedLevel2
	Name string `json:"name_level1"`
}
type TestTraverse struct {
	*TestEmbeddedLevel1
	// TestEmbeddedLevel1
	Name   string `json:"name"`
	Age    int    `json:"age"`
	gender string `json:"gender"`
}

func TestGetReflectValue(t *testing.T) {
	tt := TestTraverse{}
	_, err := getValueOfPointer(tt)
	siutils.AssertNotNilFail(t, err)

	_, err = getValueOfPointer(&tt)
	siutils.AssertNilFail(t, err)

	ttSlice := []TestTraverse{}
	_, err = getValueOfPointer(&ttSlice)
	siutils.AssertNilFail(t, err)

}

func TestTraverseFields(t *testing.T) {
	tt := []TestTraverse{}
	rvSlice := reflect.Indirect(reflect.ValueOf(&tt))
	rvTypeElem := rvSlice.Type().Elem() // type of the slice element
	rvTypeElemKind := rvTypeElem.Kind()
	var elem reflect.Value // slice's element
	if rvTypeElemKind == reflect.Pointer {
		elem = reflect.New(rvTypeElem.Elem())
	} else if rvTypeElemKind == reflect.Struct {
		elem = reflect.New(rvTypeElem).Elem()
	}

	var traversedFields []traversedField
	var fieldsToInitialize [][]int
	traverseFields(traversedField{elem, []int{}}, &traversedFields, &fieldsToInitialize)

	// fmt.Println(traversedFields)
	// for _, v := range traversedFields {
	// 	// fmt.Println(elem.FieldByIndex(v.indices).Type())
	// 	field := elem.Type().FieldByIndex(v.indices)
	// 	name, _ := findFieldNameByTag("json", field.Tag)
	// 	fmt.Println(field.Type, name)
	// }

	tagMap := makeNameMap(elem, "json", traversedFields)
	fmt.Println(tagMap)
}
