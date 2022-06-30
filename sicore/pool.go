package sicore

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
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

func getRowScanner(opts ...RowScannerOption) *RowScanner {
	rs := _rowScannerPool.Get().(*RowScanner)
	rs.Reset(opts...)
	return rs
}
func putRowScanner(rs *RowScanner) {
	rs.Reset()
	_rowScannerPool.Put(rs)
}

// GetRowScanner retrieves RowScanner from a pool or creates a new.
func GetRowScanner(opts ...RowScannerOption) *RowScanner {
	return getRowScanner(opts...)
}

// PutRowScanner puts RowScanner back to the pool.
func PutRowScanner(rs *RowScanner) {
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

// GetReader retrieves Reader from a pool or creates a new.
func GetReader(r io.Reader, opt ...ReaderOption) *Reader {
	return getReader(r, opt...)
}

// PutReader puts Reader back to the pool.
func PutReader(r *Reader) {
	putReader(r)
}

var (
	_writerPool = sync.Pool{}
)

func getWriter(w io.Writer, opt ...WriterOption) *Writer {
	g := _writerPool.Get()

	var wr *Writer
	if g == nil {
		wr = newWriter(w, opt...)
	} else {
		wr = g.(*Writer)
		wr.Reset(w, opt...)
	}
	return wr
}

func putWriter(w *Writer) {
	w.Reset(nil)
	_writerPool.Put(w)
}

// GetWriter retrieves Writer from a pool or creates a new.
func GetWriter(w io.Writer, opt ...WriterOption) *Writer {
	return getWriter(w, opt...)
}

func GetWriterAndBuffer(opt ...WriterOption) (*Writer, *bytes.Buffer) {
	buf := GetBytesBuffer(nil)
	wr := GetWriter(buf, opt...)
	return wr, buf
}

// PutWriter puts Writer back to the pool.
func PutWriter(w *Writer) {
	putWriter(w)
}

func PutWriterAndBuffer(w *Writer, buf *bytes.Buffer) {
	PutWriter(w)
	PutBytesBuffer(buf)
}

var (
	_readwriterPool = sync.Pool{}
)

func getReadWriter(r io.Reader, w io.Writer) *ReadWriter {
	g := _readwriterPool.Get()
	if g == nil {
		rd := GetReader(r)
		wr := GetWriter(w)
		return newReadWriter(rd, wr)
	}
	rw := g.(*ReadWriter)
	rw.Reader.Reset(r)
	rw.Writer.Reset(w)

	return rw
}

func putReadWriter(rw *ReadWriter) {
	rw.Reader.Reset(nil)
	rw.Writer.Reset(nil)
	_readwriterPool.Put(rw)
}

// GetReadWriter retrieves ReadWriter from a pool or creates a new.
func GetReadWriter(r io.Reader, w io.Writer) *ReadWriter {
	return getReadWriter(r, w)
}

func GetReadWriterWithReadWriter(rw io.ReadWriter) *ReadWriter {
	return getReadWriter(rw, rw)
}

// GetReadWriterWithOptions retrieves ReadWriter from a pool or creates a new with Reader and Writer options.
// func GetReadWriterWithOptions(r io.Reader, ro []ReaderOption, w io.Writer, wo []WriterOption) *ReadWriter {
// 	return getReadWriter(r, ro, w, wo)
// }

// PutReadWriter puts ReadWriter back to the pool.
func PutReadWriter(rw *ReadWriter) {
	putReadWriter(rw)
}

var (
	// bytes.Reader pool
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

var (
	// bytes.Buffer pool
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

var (
	// MapSlice pool
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

// HmacSha256HashPool wraps Get and Put methods.
type HmacSha256HashPool interface {
	Get() interface{}
	Put(v interface{})
}

var (
	_hmacSha256HashMap sync.Map
)

func getHmacSha256Pool(secret string) HmacSha256HashPool {
	p, _ := _hmacSha256HashMap.LoadOrStore(secret, &sync.Pool{
		New: func() interface{} {
			return hmac.New(sha256.New, []byte(secret))
		},
	})

	return p.(HmacSha256HashPool)
}

// GetHmacSha256Hash retrieve a Hash with secret from a pool or create a new.
func GetHmacSha256Hash(secret string) hash.Hash {
	return getHmacSha256Pool(secret).Get().(hash.Hash)
}

// PutHmacSha256Hash puts Hash with secret back into a pool.
func PutHmacSha256Hash(secret string, h hash.Hash) {
	h.Reset()
	getHmacSha256Pool(secret).Put(h)
}
