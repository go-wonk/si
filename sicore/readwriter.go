package sicore

import (
	"bufio"
	"errors"
	"io"
)

const defaultBufferSize = 512

// Flusher interface has Flush method to check if a writer has a flush method like bufio.Writer.
// json.Encoder doesn't flush after write.
type Flusher interface {
	Flush() error
}

// Lengther wraps Len medhod
type Lengther interface {
	Len() int
}

type WriterResetter interface {
	Reset(w io.Writer)
}

type ReaderResetter interface {
	Reset(r io.Reader)
}

// Reader is a wrapper of buifio.Reader.
type Reader struct {
	r   io.Reader
	br  *bufio.Reader
	dec Decoder
	chk EofChecker

	bufAll []byte
}

func newReader(r io.Reader, opt ...ReaderOption) *Reader {
	var br *bufio.Reader
	var ok bool
	if br, ok = r.(*bufio.Reader); !ok {
		br = bufio.NewReader(r)
	}

	rd := &Reader{r: r, br: br, bufAll: make([]byte, 0, defaultBufferSize)}
	rd.ApplyOptions(opt...)
	return rd
}

func (rd *Reader) ApplyOptions(opts ...ReaderOption) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(rd)
	}

	// always set DefaultEofChecker if `rd.chk` is not set
	if rd.r != nil && rd.chk == nil {
		rd.chk = DefaultEofChecker
	}
	if rd.r != nil && rd.dec == nil {
		rd.dec = &DefaultDecoder{rd}
	}
}

// SetEofChecker sets EofChecker to underlying Reader.
func (rd *Reader) SetEofChecker(chk EofChecker) {
	rd.chk = chk
}

func (rd *Reader) SetDecoder(dec Decoder) {
	rd.dec = dec
}

// Reset resets underlying Reader with r and opt.
func (rd *Reader) Reset(r io.Reader, opt ...ReaderOption) {
	rd.bufAll = rd.bufAll[:0]
	rd.r = r
	rd.br.Reset(r)

	if rs, ok := rd.dec.(ReaderResetter); ok {
		rs.Reset(rd)
	} else {
		rd.dec = nil
	}

	rd.chk = nil

	rd.ApplyOptions(opt...)
}

// Read reads the data of underlying Reader(rd.br) into p.
func (rd *Reader) Read(p []byte) (n int, err error) {
	n, err = rd.br.Read(p)
	return
}

// ReadAll reads all data from underlying Reader(rd.br) and returns it.
func (rd *Reader) ReadAll() ([]byte, error) {
	// return readAll(rd.br, rd.chk)
	return rd.readAll()
}

func (rd *Reader) readAll() ([]byte, error) {
	rd.bufAll = rd.bufAll[:0]
	for {
		if len(rd.bufAll) == cap(rd.bufAll) {
			if err := growCap(&rd.bufAll, defaultBufferSize); err != nil {
				return nil, err
			}
		}

		n, err := rd.br.Read(rd.bufAll[len(rd.bufAll):cap(rd.bufAll)])
		rd.bufAll = rd.bufAll[:len(rd.bufAll)+n]

		ended, err := rd.chk.Check(rd.bufAll, err)
		if err != nil {
			rb := make([]byte, len(rd.bufAll))
			n := copy(rb, rd.bufAll[:len(rd.bufAll)])
			return rb[:n], err
		}
		if ended {
			rb := make([]byte, len(rd.bufAll))
			n := copy(rb, rd.bufAll[:len(rd.bufAll)])
			return rb[:n], nil
		}
	}
}

// readAll reads all data from r and returns it
func readAll(r io.Reader, chk EofChecker) ([]byte, error) {

	b := make([]byte, 0, defaultBufferSize)
	for {
		if len(b) == cap(b) {
			if err := growCap(&b, defaultBufferSize); err != nil {
				return nil, err
			}
		}

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		ended, err := chk.Check(b, err)
		if err != nil {
			return b, err
		}
		if ended {
			return b, nil
		}
	}
}

var ErrNoDecoder = errors.New("no decoder was provided")

// Decode decodes data from underlying Reader(rd.br) and saves it to the value pointed by v.
func (rd *Reader) Decode(v any) error {
	if rd.dec == nil {
		return ErrNoDecoder
	}
	return rd.dec.Decode(v)
}

// Peek returns next n bytes of underlying Reader(rd.br) without advancing the Reader.
func (rd *Reader) Peek(n int) ([]byte, error) {
	return rd.br.Peek(n)
}

// Buffered returns the number of bytes that can be read from the buffer of underlying Reader(rd.br).
func (rd *Reader) Buffered() int {
	return rd.br.Buffered()
}

// Size returns size of buffer of underlying Reader(rd.br).
// eg. len(buf)
func (rd *Reader) Size() int {
	return rd.br.Size()
}

// Len returns the length of underlying Reader(rd.r) if rd.r implements Len() method.
// It returns 0 otherwise.
// This is different from Buffered. Buffered returns the length of temporary data in a buffer.
func (rd *Reader) Len() int {
	if l, ok := rd.r.(Lengther); ok {
		return l.Len()
	}
	return 0
}

// WriteTo writes data of underlying Reader(rd.br) into w.
func (rd *Reader) WriteTo(w io.Writer) (n int64, err error) {
	return rd.br.WriteTo(w)
}

// Writer is a wrapper of bufio.Writer with Encoder.
type Writer struct {
	w   io.Writer
	bw  *bufio.Writer
	enc Encoder
}

func newWriter(w io.Writer, opt ...WriterOption) *Writer {
	var bw *bufio.Writer
	var ok bool
	if bw, ok = w.(*bufio.Writer); !ok {
		bw = bufio.NewWriter(w)
	}
	wr := &Writer{w: w, bw: bw}
	wr.ApplyOptions(opt...)
	return wr
}

func (wr *Writer) ApplyOptions(opts ...WriterOption) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(wr)
	}
	if wr.w != nil && wr.enc == nil {
		wr.enc = &DefaultEncoder{wr}
	}
}

func (wr *Writer) SetEncoder(enc Encoder) {
	wr.enc = enc
}

func (wr *Writer) Available() int {
	return wr.bw.Available()
}

// Buffered returns the number of bytes that have been written to the buffer.
func (wr *Writer) Buffered() int {
	return wr.bw.Buffered()
}

// Flush writes the data left in buffer to the underlying Writer(wr.bw).
func (wr *Writer) Flush() error {
	return wr.bw.Flush()
}

// ReadFrom reads from r into underlying Writer(wr.bw).
func (wr *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = wr.bw.ReadFrom(r)
	return
}

// Reset resets underlying Writer(wr) with w and opt.
func (wr *Writer) Reset(w io.Writer, opt ...WriterOption) {
	wr.w = w
	wr.bw.Reset(w)

	if rs, ok := wr.enc.(WriterResetter); ok {
		rs.Reset(wr)
	} else {
		wr.enc = nil
	}
	wr.ApplyOptions(opt...)
}

func (wr *Writer) Size() int {
	return wr.bw.Size()
}

// Write writes p into underlying Writer(wr.bw).
func (wr *Writer) Write(p []byte) (n int, err error) {
	n, err = wr.bw.Write(p)
	return
}

func (wr *Writer) WriteByte(c byte) error {
	return wr.bw.WriteByte(c)
}

func (wr *Writer) WriteRune(r rune) (size int, err error) {
	return wr.bw.WriteRune(r)
}

func (wr *Writer) WriteString(s string) (n int, err error) {
	return wr.bw.WriteString(s)
}

// WriteFlush writes p to underlying Writer followed by Flush.
func (wr *Writer) WriteFlush(p []byte) (n int, err error) {
	n, err = wr.Write(p)
	if err != nil {
		return
	}
	if err = wr.Flush(); err != nil {
		n = 0
		return
	}
	return
}

var ErrNoEncoder = errors.New("no encoder was provided")

// Encode writes encoded data into underlying Writer(wr.bw).
func (wr *Writer) Encode(p any) (err error) {
	if wr.enc == nil {
		return ErrNoEncoder
	}
	err = wr.enc.Encode(p)
	return
}

// EncodeFlush writes encoded data into underlying Writer.
// It flushes any data remaining in the buffer right away.
func (wr *Writer) EncodeFlush(p any) (err error) {
	if err = wr.Encode(p); err != nil {
		return
	}
	err = wr.Flush()
	return
}

// // ReadWriter uses bufio package to read and write more efficiently.
// // It is designed to read/write data from/to a storage that implements ReadWriter interface.
// // `validator` determines when to finish reading, and defaultValidator is to finish when io.EOF is met.
type ReadWriter struct {
	*Reader
	*Writer
}

func newReadWriter(r *Reader, w *Writer) *ReadWriter {
	return &ReadWriter{r, w}
}

func (rw *ReadWriter) RLen() int {
	return rw.Reader.Len()
}

func (rw *ReadWriter) RBuffered() int {
	return rw.Reader.Buffered()
}

func (rw *ReadWriter) WBuffered() int {
	return rw.Writer.Buffered()
}

func (rw *ReadWriter) Request(p []byte) ([]byte, error) {
	_, err := rw.WriteFlush(p)
	if err != nil {
		return nil, err
	}
	b, err := rw.ReadAll()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (rw *ReadWriter) RequestEncoded(v any) ([]byte, error) {
	if err := rw.EncodeFlush(v); err != nil {
		return nil, err
	}
	b, err := rw.ReadAll()
	if err != nil {
		return nil, err
	}
	return b, nil
}
