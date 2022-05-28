package sicore

import (
	"io"
	"sync"
)

// var (
// 	_bioReaderPool = sync.Pool{
// 		New: func() interface{} {
// 			return new(bufio.Reader)
// 		},
// 	}

// 	_bioWriterPool = sync.Pool{
// 		New: func() interface{} {
// 			return new(bufio.Writer)
// 		},
// 	}
// )

// func getBufioReader(r io.Reader) *bufio.Reader {
// 	br := _bioReaderPool.Get().(*bufio.Reader)
// 	br.Reset(r)
// 	return br
// }
// func putBufioReader(br *bufio.Reader) {
// 	br.Reset(nil)
// 	_bioReaderPool.Put(br)
// }

// func GetBufioReader(r io.Reader) *bufio.Reader {
// 	return getBufioReader(r)
// }
// func PutBufioReader(br *bufio.Reader) {
// 	putBufioReader(br)
// }

// func getBufioWriter(w io.Writer) *bufio.Writer {
// 	bw := _bioWriterPool.Get().(*bufio.Writer)
// 	bw.Reset(w)
// 	return bw
// }
// func putBufioWriter(bw *bufio.Writer) {
// 	bw.Reset(nil)
// 	_bioWriterPool.Put(bw)
// }

// func GetBufioWriter(w io.Writer) *bufio.Writer {
// 	return getBufioWriter(w)
// }
// func PutBufioWriter(bw *bufio.Writer) {
// 	putBufioWriter(bw)
// }

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

func getWriter(w io.Writer, enc EncoderSetter) *Writer {
	g := _writerPool.Get()
	if g == nil {
		return newWriter(w, enc)
	}
	wr := g.(*Writer)
	wr.Reset(w, enc)
	return wr
}

func putWriter(w *Writer) {
	w.Reset(nil, nil)
	_writerPool.Put(w)
}

func GetWriter(w io.Writer) *Writer {
	return getWriter(w, nil)
}
func GetWriterWithEncoder(w io.Writer, enc EncoderSetter) *Writer {
	return getWriter(w, enc)
}
func PutWriter(w *Writer) {
	putWriter(w)
}
