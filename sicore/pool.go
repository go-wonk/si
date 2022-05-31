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

func getRowScanner(useSqlNullType bool) *rowScanner {
	rs := _rowScannerPool.Get().(*rowScanner)
	rs.Reset(useSqlNullType)
	return rs
}
func putRowScanner(rs *rowScanner) {
	rs.Reset(defaultUseSqlNullType)
	_rowScannerPool.Put(rs)
}

func GetRowScanner() *rowScanner {
	return getRowScanner(defaultUseSqlNullType)
}
func PutRowScanner(rs *rowScanner) {
	putRowScanner(rs)
}

var (
	_readerPool = sync.Pool{}
)

func getReader(r io.Reader, opt ...ReaderOption) *Reader {
	g := _readerPool.Get()
	if g == nil {
		return newReader(r, opt...)
	}
	rd := g.(*Reader)
	rd.Reset(r, opt...)
	return rd
}
func putReader(r *Reader) {
	r.Reset(nil)
	_readerPool.Put(r)
}

func GetReader(r io.Reader, opt ...ReaderOption) *Reader {
	return getReader(r, opt...)
}
func PutReader(r *Reader) {
	putReader(r)
}

var (
	_writerPool = sync.Pool{}
)

func getWriter(w io.Writer, opt ...WriterOption) *Writer {
	g := _writerPool.Get()
	if g == nil {
		return newWriter(w, opt...)
	}
	wr := g.(*Writer)
	wr.Reset(w, opt...)
	return wr
}

func putWriter(w *Writer) {
	w.Reset(nil)
	_writerPool.Put(w)
}

func GetWriter(w io.Writer, opt ...WriterOption) *Writer {
	return getWriter(w, opt...)
}
func PutWriter(w *Writer) {
	putWriter(w)
}

var (
	_readwriterPool = sync.Pool{}
)

func getReadWriter(r io.Reader, ro []ReaderOption, w io.Writer, wo []WriterOption) *ReadWriter {
	g := _readwriterPool.Get()
	if g == nil {
		rd := GetReader(r, ro...)
		wr := GetWriter(w, wo...)
		return newReadWriter(rd, wr)
	}
	rw := g.(*ReadWriter)
	rw.Reader.Reset(r, ro...)
	rw.Writer.Reset(w, wo...)

	return rw
}

func putReadWriter(rw *ReadWriter) {
	rw.Reader.Reset(nil)
	rw.Writer.Reset(nil)
	_readwriterPool.Put(rw)
}

func GetReadWriter(r io.Reader, ro []ReaderOption, w io.Writer, wo []WriterOption) *ReadWriter {
	return getReadWriter(r, ro, w, wo)
}

func PutReadWriter(rw *ReadWriter) {
	putReadWriter(rw)
}
