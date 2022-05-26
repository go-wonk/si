package sicore

import "database/sql"

type SqlColType uint8

const (
	SqlColTypeBool SqlColType = iota
	SqlColTypeByte
	SqlColTypeBytes
	SqlColTypeString
	SqlColTypeInt
	SqlColTypeInt8
	SqlColTypeInt16
	SqlColTypeInt32
	SqlColTypeInt64
	SqlColTypeUint
	SqlColTypeUint8
	SqlColTypeUint16
	SqlColTypeUint32
	SqlColTypeUint64
	SqlColTypeFloat32
	SqlColTypeFloat64
	SqlColTypeTime
	SqlColTypeints
	SqlColTypeints8
	SqlColTypeints16
	SqlColTypeints32
	SqlColTypeints64
	SqlColTypeUints
	SqlColTypeUints8
	SqlColTypeUints16
	SqlColTypeUints32
	SqlColTypeUints64
)

var (
	// stringTypeValue  string
	// bytesTypeValue   []byte
	// intTypeValue     int
	// int8TypeValue    int8
	// int16TypeValue   int16
	// int32TypeValue   int32
	// int64TypeValue   int64
	// uintTypeValue    uint
	// uint8TypeValue   uint8
	// uint16TypeValue  uint16
	// uint32TypeValue  uint32
	// uint64TypeValue  uint64
	// boolTypeValue    bool
	// float32TypeValue float32
	// float64TypeValue float64
	// timeTypeValue    time.Time
	// byteTypeValue    byte
	intsTypeValue   []int
	ints8TypeValue  []int8
	ints16TypeValue []int16
	ints32TypeValue []int32
	ints64TypeValue []int64

	uintsTypeValue   []uint
	uints8TypeValue  []uint8
	uints16TypeValue []uint16
	uints32TypeValue []uint32
	uints64TypeValue []uint64

	sqlNullBoolTypeValue    sql.NullBool
	sqlNullByteTypeValue    sql.NullByte
	sqlBytesTypeValue       sql.RawBytes
	sqlNullStringTypeValue  sql.NullString
	sqlNullFloat32TypeValue sql.NullFloat64
	sqlNullFloat64TypeValue sql.NullFloat64
	sqlNullIntTypeValue     sql.NullInt64
	sqlNullInt8TypeValue    sql.NullInt16
	sqlNullInt16TypeValue   sql.NullInt16
	sqlNullInt32TypeValue   sql.NullInt32
	sqlNullInt64TypeValue   sql.NullInt64
	sqlNullUintTypeValue    sql.NullInt64
	sqlNullUint8TypeValue   sql.NullInt16
	sqlNullUint16TypeValue  sql.NullInt16
	sqlNullUint32TypeValue  sql.NullInt32
	sqlNullUint64TypeValue  sql.NullInt64
	sqlNullTimeTypeValue    sql.NullTime
)

type SqlColumn struct {
	Name string
	Type SqlColType
}

func (sc SqlColumn) SetType(rs *rowScanner) {
	switch sc.Type {
	case SqlColTypeBool:
		rs.SetSqlColumn(sc.Name, sqlNullBoolTypeValue)
	case SqlColTypeByte:
		rs.SetSqlColumn(sc.Name, sqlNullByteTypeValue)
	case SqlColTypeBytes:
		rs.SetSqlColumn(sc.Name, sqlBytesTypeValue)
	case SqlColTypeString:
		rs.SetSqlColumn(sc.Name, sqlNullStringTypeValue)
	case SqlColTypeInt:
		rs.SetSqlColumn(sc.Name, sqlNullIntTypeValue)
	case SqlColTypeInt8:
		rs.SetSqlColumn(sc.Name, sqlNullInt8TypeValue)
	case SqlColTypeInt16:
		rs.SetSqlColumn(sc.Name, sqlNullInt16TypeValue)
	case SqlColTypeInt32:
		rs.SetSqlColumn(sc.Name, sqlNullInt32TypeValue)
	case SqlColTypeInt64:
		rs.SetSqlColumn(sc.Name, sqlNullInt64TypeValue)
	case SqlColTypeUint:
		rs.SetSqlColumn(sc.Name, sqlNullUintTypeValue)
	case SqlColTypeUint8:
		rs.SetSqlColumn(sc.Name, sqlNullUint8TypeValue)
	case SqlColTypeUint16:
		rs.SetSqlColumn(sc.Name, sqlNullUint16TypeValue)
	case SqlColTypeUint32:
		rs.SetSqlColumn(sc.Name, sqlNullUint32TypeValue)
	case SqlColTypeUint64:
		rs.SetSqlColumn(sc.Name, sqlNullUint64TypeValue)
	case SqlColTypeFloat32:
		rs.SetSqlColumn(sc.Name, sqlNullFloat32TypeValue)
	case SqlColTypeFloat64:
		rs.SetSqlColumn(sc.Name, sqlNullFloat64TypeValue)
	case SqlColTypeTime:
		rs.SetSqlColumn(sc.Name, sqlNullTimeTypeValue)
	case SqlColTypeints:
		rs.SetSqlColumn(sc.Name, intsTypeValue)
	case SqlColTypeints8:
		rs.SetSqlColumn(sc.Name, ints8TypeValue)
	case SqlColTypeints16:
		rs.SetSqlColumn(sc.Name, ints16TypeValue)
	case SqlColTypeints32:
		rs.SetSqlColumn(sc.Name, ints32TypeValue)
	case SqlColTypeints64:
		rs.SetSqlColumn(sc.Name, ints64TypeValue)
	case SqlColTypeUints:
		rs.SetSqlColumn(sc.Name, uintsTypeValue)
	case SqlColTypeUints8:
		rs.SetSqlColumn(sc.Name, uints8TypeValue)
	case SqlColTypeUints16:
		rs.SetSqlColumn(sc.Name, uints16TypeValue)
	case SqlColTypeUints32:
		rs.SetSqlColumn(sc.Name, uints32TypeValue)
	case SqlColTypeUints64:
		rs.SetSqlColumn(sc.Name, uints64TypeValue)
	}
}
