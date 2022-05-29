package si

import "github.com/go-wonk/si/sicore"

func SqlFloat64(columnName string) sicore.SqlColumn {
	return sicore.SqlColumn{Name: columnName, Type: sicore.SqlColTypeFloat64}
}

func SqlString(columnName string) sicore.SqlColumn {
	return sicore.SqlColumn{Name: columnName, Type: sicore.SqlColTypeString}
}

func SqlUints8(columnName string) sicore.SqlColumn {
	return sicore.SqlColumn{Name: columnName, Type: sicore.SqlColTypeUints8}
}
