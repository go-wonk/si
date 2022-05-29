package sicore

import (
	"io"
	"sync"
)

var (
	_rowScannerPool = sync.Pool{
		New: func() interface{} {
			return newRowScanner()
		},
	}
)

func getRowScanner(sqlCol map[string]any, useSqlNullType bool) *rowScanner {
	rs := _rowScannerPool.Get().(*rowScanner)
	rs.Reset(sqlCol, useSqlNullType)
	return rs
}
func putRowScanner(rs *rowScanner) {
	rs.Reset(nil, defaultUseSqlNullType)
	_rowScannerPool.Put(rs)
}

func GetRowScanner() *rowScanner {
	return getRowScanner(make(map[string]any), defaultUseSqlNullType)
}
func PutRowScanner(rs *rowScanner) {
	putRowScanner(rs)
}

var (
	_readerPool = sync.Pool{}
)

func getReader(r io.Reader, val ReadValidator) *Reader {
	g := _readerPool.Get()
	if g == nil {
		return newReader(r, val)
	}
	rd := g.(*Reader)
	rd.Reset(r, val)
	return rd
}
func putReader(r *Reader) {
	r.Reset(nil, nil)
	_readerPool.Put(r)
}

func GetReader(r io.Reader) *Reader {
	return getReader(r, DefaultValidator())
}
func GetReaderWithValidator(r io.Reader, val ReadValidator) *Reader {
	return getReader(r, val)
}
func PutReader(r *Reader) {
	putReader(r)
}

var (
	_writerPool = sync.Pool{}
)

func getWriter(w io.Writer, opt ...Option) *Writer {
	g := _writerPool.Get()
	if g == nil {
		return newWriter(w, opt...)
	}
	wr := g.(*Writer)
	wr.Reset(w, opt...)
	return wr
}

func putWriter(w *Writer) {
	w.Reset(nil, nil)
	_writerPool.Put(w)
}

func GetWriter(w io.Writer, opt ...Option) *Writer {
	return getWriter(w, opt...)
}
func PutWriter(w *Writer) {
	putWriter(w)
}
