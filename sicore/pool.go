package sicore

import (
	"bufio"
	"io"
	"sync"
)

var (
	_bioReaderPool = sync.Pool{
		New: func() interface{} {
			return new(bufio.Reader)
		},
	}

	_bioWriterPool = sync.Pool{
		New: func() interface{} {
			return new(bufio.Writer)
		},
	}
)

func getBufioReader(r io.Reader) *bufio.Reader {
	br := _bioReaderPool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}
func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	_bioReaderPool.Put(br)
}

func GetBufioReader(r io.Reader) *bufio.Reader {
	return getBufioReader(r)
}
func PutBufioReader(br *bufio.Reader) {
	putBufioReader(br)
}

func getBufioWriter(w io.Writer) *bufio.Writer {
	bw := _bioWriterPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}
func putBufioWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	_bioWriterPool.Put(bw)
}

func GetBufioWriter(w io.Writer) *bufio.Writer {
	return getBufioWriter(w)
}
func PutBufioWriter(bw *bufio.Writer) {
	putBufioWriter(bw)
}

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
	_readerPool = sync.Pool{
		New: func() interface{} {
			return newReader()
		},
	}
)

func getReader(r io.Reader, bufferSize int) *Reader {
	rd := _readerPool.Get().(*Reader)
	rd.Reset(r, bufferSize, defaultValidate())
	return rd
}
func putReader(r *Reader) {
	r.Reset(nil, defaultBufferSize, nil)
	_rowScannerPool.Put(r)
}

func GetReader(r io.Reader) *Reader {
	return GetReaderSize(r, defaultBufferSize)
}
func GetReaderSize(r io.Reader, bufferSize int) *Reader {
	return getReader(r, bufferSize)
}
func PutReader(r *Reader) {
	putReader(r)
}
