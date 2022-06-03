package sicore

import (
	"bytes"
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
	// smallBufferSize  = 128
	// mediumBufferSize = 1024

	_readerPool = sync.Pool{}
	// _readerPoolSmall  = sync.Pool{}
	// _readerPoolMedium = sync.Pool{}
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

func GetReadWriter(r io.Reader, w io.Writer) *ReadWriter {
	return getReadWriter(r, nil, w, nil)
}

func GetReadWriterWithOptions(r io.Reader, ro []ReaderOption, w io.Writer, wo []WriterOption) *ReadWriter {
	return getReadWriter(r, ro, w, wo)
}

func PutReadWriter(rw *ReadWriter) {
	putReadWriter(rw)
}

// bytes.Reader pool
var (
	_bytesReaderPool = sync.Pool{}
)

func getBytesReader(b []byte) *bytes.Reader {
	g := _bytesReaderPool.Get()
	if g == nil {
		return bytes.NewReader(b)
	}
	br := g.(*bytes.Reader)
	br.Reset(b)
	return br
}

func putBytesReader(r *bytes.Reader) {
	_bytesReaderPool.Put(r)
}

func GetBytesReader(b []byte) *bytes.Reader {
	return getBytesReader(b)
}

func PutBytesReader(r *bytes.Reader) {
	putBytesReader(r)
}

// bytes.Buffer pool
var (
	_bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 512))
		},
	}
)

func getBytesBuffer(b []byte) *bytes.Buffer {
	bb := _bytesBufferPool.Get().(*bytes.Buffer)
	bb.Reset()
	if len(b) > 0 {
		_, err := bb.Write(b)
		if err != nil {
			return bytes.NewBuffer(b)
		}
	}
	return bb
}

func putBytesBuffer(r *bytes.Buffer) {
	_bytesBufferPool.Put(r)
}

func GetBytesBuffer(b []byte) *bytes.Buffer {
	return getBytesBuffer(b)
}

func PutBytesBuffer(b *bytes.Buffer) {
	putBytesBuffer(b)
}

// MapSlice pool
var (
	_msPool = sync.Pool{
		New: func() interface{} {
			ms := make([]map[string]interface{}, 0, 100)
			return ms
		},
	}
)

func getMapSlice() []map[string]interface{} {
	ms := _msPool.Get().([]map[string]interface{})
	return ms
}

func putMapSlice(ms []map[string]interface{}) {
	for i := range ms {
		for k := range ms[i] {
			// ms[i][k] = nil //
			delete(ms[i], k)
		}
	}
	ms = ms[:0]
	_msPool.Put(ms)
}

// const maxInt = int(^uint(0) >> 1)

// var ErrTooLarge = errors.New("buf too large")

func growMapSlice(ms *[]map[string]interface{}, s int) (int, error) {
	c := cap(*ms)
	l := len(*ms)
	a := c - l // available
	if s <= a {
		*ms = (*ms)[:l+s]
		return l, nil
	}

	if l+s <= c {
		// if needed length is lte c
		return l, nil
	}

	if c > maxInt-c-s {
		// too large
		return l, ErrTooLarge
	}

	newBuf := make([]map[string]interface{}, c*2+s)
	copy(newBuf, (*ms)[0:])
	*ms = newBuf[:l+s]
	return l, nil
}

func makeMapIfNil(m *map[string]interface{}) {
	if *m == nil {
		*m = make(map[string]interface{})
		return
	}
}
