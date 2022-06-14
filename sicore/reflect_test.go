package sicore

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
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
	_, err := valueOfAnyPtr(tt)
	siutils.AssertNotNilFail(t, err)

	_, err = valueOfAnyPtr(&tt)
	siutils.AssertNilFail(t, err)

	ttSlice := []TestTraverse{}
	_, err = valueOfAnyPtr(&ttSlice)
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

func TestStructReflectType(t *testing.T) {
	s := testmodels.Student{}
	columns := []string{"id", "name", "email_address", "borrowed"}

	rv := reflect.ValueOf(&s)
	if rv.Kind() != reflect.Pointer {
		t.FailNow()
	}
	fmt.Println(rv) // {"id":0,"email_address":"","name":"","borrowed":false}

	rve := rv.Elem()
	fmt.Println(rve) // {0   false null}

	rvType := rv.Type() // just checking its type
	fmt.Println(rvType) // *testmodels.Student

	rveType := rve.Type() // just checking its type
	fmt.Println(rveType)  // testmodels.Student

	// initializeFields(rv) // this panics because rv is a pointer

	// initializeFields(rve) // initializes any nil pointer struct fields
	// fmt.Println(rve)      // {0   false {"book_id":0}}

	var traversedFields []traversedField
	var fieldsToInitialize [][]int
	traverseFields(traversedField{rve, []int{}}, &traversedFields, &fieldsToInitialize)
	fmt.Println(rve)                // {0   false {"book_id":0}}
	fmt.Println(traversedFields)    // [{{0x1005d43a0 0x1400002d180 386} [0]} {{0x1005d4d20 0x1400002d188 408} [1]} {{0x1005d4d20 0x1400002d198 408} [2]} {{0x1005d2ce0 0x1400002d1a8 385} [3]} {{0x1005d43a0 0x140000190d0 386} [4 0]}]
	fmt.Println(fieldsToInitialize) // [[4]]

	tagNameMap := makeNameMap(rve, "json", traversedFields)
	fmt.Println(tagNameMap) // map[book_id:[4 0] borrowed:[3] email_address:[1] id:[0] name:[2]]

	scannedRow, err := buildDestinations(columns, tagNameMap, rve)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(scannedRow...)
}

func TestPrimaryReflectType(t *testing.T) {
	var id int
	// columns := []string{"id"}

	rv := reflect.ValueOf(&id)
	if rv.Kind() != reflect.Pointer {
		t.FailNow()
	}
	fmt.Println(rv) // 0x14000018ea0

	rve := rv.Elem()
	fmt.Println(rve) // 0

	rvType := rv.Type() // just checking its type
	fmt.Println(rvType) // *int

	rveType := rve.Type() // just checking its type
	fmt.Println(rveType)  // int

	switch rve.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Println("is Int")
	case reflect.Struct:
		fmt.Println("is struct")
	default:
		fmt.Println(rve.Type().Kind())
	}

	// initializeFields(rv) // this panics because rv is a pointer

	// initializeFields(rve) // initializes any nil pointer struct fields
	// fmt.Println(rve)      // {0   false {"book_id":0}}

	// var traversedFields []traversedField
	// var fieldsToInitialize [][]int
	// traverseFields(traversedField{rve, []int{}}, &traversedFields, &fieldsToInitialize)
	// fmt.Println(rve)                // {0   false {"book_id":0}}
	// fmt.Println(traversedFields)    // [{{0x1005d43a0 0x1400002d180 386} [0]} {{0x1005d4d20 0x1400002d188 408} [1]} {{0x1005d4d20 0x1400002d198 408} [2]} {{0x1005d2ce0 0x1400002d1a8 385} [3]} {{0x1005d43a0 0x140000190d0 386} [4 0]}]
	// fmt.Println(fieldsToInitialize) // [[4]]

	// tagNameMap := makeNameMap(rve, "json", traversedFields)
	// fmt.Println(tagNameMap) // map[book_id:[4 0] borrowed:[3] email_address:[1] id:[0] name:[2]]

	// scannedRow, err := buildDestinations(columns, tagNameMap, rve)
	// if err != nil {
	// 	t.FailNow()
	// }
	// fmt.Println(scannedRow...)
}
