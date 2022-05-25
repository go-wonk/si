package sicore

import (
	"fmt"
	"testing"
)

func TestSql_Type(t *testing.T) {
	// fmt.Printf("%p, %T\n", refTypeOfBytes, refTypeOfBytes)
	// fmt.Printf("%p, %T\n", refTypeOfBytes, refTypeOfBytes)

	// fmt.Printf("%p, %T\n", reflect.TypeOf([]byte("")), reflect.TypeOf([]byte("")))

	// fmt.Printf("%p, %T\n", reflect.PtrTo(refTypeOfBytes), reflect.PtrTo(refTypeOfBytes))
	// fmt.Printf("%v, %T\n", reflect.New(reflect.PtrTo(refTypeOfBytes)), reflect.New(reflect.PtrTo(refTypeOfBytes)))
	// fmt.Printf("%p, %T\n", reflect.PtrTo(refTypeOfBytes), reflect.PtrTo(refTypeOfBytes))
	// fmt.Printf("%v, %T\n", reflect.New(reflect.PtrTo(refTypeOfBytes)), reflect.New(reflect.PtrTo(refTypeOfBytes)))

	tt := []uint8{57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57}
	aa := string(tt)
	fmt.Printf("%v\n", aa)

	var bb complex128 = 999999999999999999
	fmt.Printf("%v\n", bb)
}
