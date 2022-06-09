package testmodels

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

/*
uint8  : 0 to 255
uint16 : 0 to 65535
uint32 : 0 to 4294967295
uint64 : 0 to 18446744073709551615
int8   : -128 to 127
int16  : -32768 to 32767
int32  : -2147483648 to 2147483647
int64  : -9223372036854775808 to 9223372036854775807

const (
        MaxFloat32             = 3.40282346638528859811704183484516925440e+38  // 2**127 * (2**24 - 1) / 2**23
        SmallestNonzeroFloat32 = 1.401298464324817070923729583289916131280e-45 // 1 / 2**(127 - 1 + 23)

        MaxFloat64             = 1.797693134862315708145274237317043567981e+308 // 2**1023 * (2**53 - 1) / 2**52
        SmallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324 // 1 / 2**(1023 - 1 + 52)
)
*/
type BasicDataType struct {
	BoolValue   bool   `json:"bool_value"`
	StringValue string `json:"string_value"`
	IntValue    int    `json:"int_value"`
	Int8Value   int8   `json:"int8_value"`
	Int16Value  int16  `json:"int16_value"`
	Int32Value  int32  `json:"int32_value"`
	Int64Value  int64  `json:"int64_value"`
	UintValue   uint   `json:"uint_value"`
	Uint8Value  uint8  `json:"uint8_value"`
	Uint16Value uint16 `json:"uint16_value"`
	Uint32Value uint32 `json:"uint32_value"`
	Uint64Value uint64 `json:"uint64_value"`
	ByteValue   byte   `json:"byte_value"`
	// RuneValue    rune    `json:"rune_value"` // not supported
	Float32Value float32 `json:"float32_value"`
	Float64Value float64 `json:"float64_value"`
	// Complex64Value  complex64  `json:"complex64_value"` // not supported
	// Complex128Value complex128 `json:"complex128_value"`// not supported
	BytesValue    []byte  `json:"bytes_value"`
	BigFloatValue []uint8 `json:"big_float_value"`

	BoolPtrValue     *bool    `json:"bool_ptr_value"`
	StringPtrValue   *string  `json:"string_ptr_value"`
	IntPtrValue      *int     `json:"int_ptr_value"`
	Int8PtrValue     *int8    `json:"int8_ptr_value"`
	Int16PtrValue    *int16   `json:"int16_ptr_value"`
	Int32PtrValue    *int32   `json:"int32_ptr_value"`
	Int64PtrValue    *int64   `json:"int64_ptr_value"`
	UintPtrValue     *uint    `json:"uint_ptr_value"`
	Uint8PtrValue    *uint8   `json:"uint8_ptr_value"`
	Uint16PtrValue   *uint16  `json:"uint16_ptr_value"`
	Uint32PtrValue   *uint32  `json:"uint32_ptr_value"`
	Uint64PtrValue   *uint64  `json:"uint64_ptr_value"`
	BytePtrValue     *byte    `json:"byte_ptr_value"`
	Float32PtrValue  *float32 `json:"float32_ptr_value"`
	Float64PtrValue  *float64 `json:"float64_ptr_value"`
	BytesPtrValue    *[]byte  `json:"bytes_ptr_value"`
	BigFloatPtrValue *[]uint8 `json:"big_float_ptr_value"`
}

func (bd *BasicDataType) String() string {
	b, err := json.Marshal(bd)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

type BasicDataTypeList []BasicDataType

func (bd *BasicDataTypeList) String() string {
	b, err := json.Marshal(bd)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

type Tabler interface {
	Query() error
}

type TableInterface interface{}

type TableString string

type ChildTable struct {
	unexportesName string `json:"unexported_child_name"`
	Name           string `json:"child_name"`
}

type Table struct {
	Nil            string         `json:"nil"`
	Int2           int            `json:"int2_"`
	Decimal        float64        `json:"decimal_"`
	Numeric        float64        `json:"numeric_"`
	Bigint         float64        `json:"bigint_"`
	CharArr        []uint8        `json:"char_arr_"`
	VarcharArr     []uint8        `json:"varchar_arr_"`
	Bytea          []byte         `json:"bytea_"`
	Time           time.Time      `json:"time_"`
	TimePtr        *time.Time     `json:"time_ptr_"`
	c1             ChildTable     `json:"c1"`
	c2             *ChildTable    `json:"c2"`
	C3             ChildTable     `json:"c3"`
	C4             *ChildTable    `json:"c4"`
	unexportesName string         `json:"unexported_name"`
	TableString    TableString    `json:"table_str"`
	TableString2   *TableString   `json:"table_str2"`
	TableInterface TableInterface `json:"table_interface"`
	Tabler         Tabler
	Any            any             `json:"any_value"`
	SqlString      sql.NullString  `json:"null_string"`
	SqlStringPtr   *sql.NullString `json:"null_string_ptr"`
}

func (t Table) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TableList []Table

func (tl TableList) String() string {
	b, _ := json.Marshal(tl)
	return string(b)
}

type TableWithNoTag struct {
	NilValue        string `json:"nil_value"`
	IntValue        int
	DecimalValue    float64
	SomeStringValue string
}

func (t *TableWithNoTag) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TableWithNoTagList []TableWithNoTag

func (t *TableWithNoTagList) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TableWithPtrElem struct {
	NilValue        string `json:"nil_value"`
	IntValue        int
	DecimalValue    float64
	SomeStringValue string
}

func (t *TableWithPtrElem) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TableWithPtrElemList []*TableWithPtrElem

func (t *TableWithPtrElemList) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}
