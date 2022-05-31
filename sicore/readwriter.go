package sicore

import (
	"bufio"
	"errors"
	"io"
)

const defaultBufferSize = 4096

// Flusher interface has Flush method to check if a writer has a flush method like bufio.Writer.
// json.Encoder doesn't flush after write.
type Flusher interface {
	Flush() error
}

// Reader
type Reader struct {
	br  *bufio.Reader
	dec Decoder
	chk EofChecker
}

func newReader(r io.Reader, opt ...ReaderOption) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		b := &Reader{br: br}
		for _, o := range opt {
			o.apply(b)
		}
		if b.chk == nil {
			b.chk = &DefaultEofChecker{}
		}
		return b
	}
	br := bufio.NewReader(r)
	b := &Reader{br: br}
	for _, o := range opt {
		o.apply(b)
	}
	if b.chk == nil {
		b.chk = &DefaultEofChecker{}
	}
	return b
}

func (rd *Reader) SetEofChecker(chk EofChecker) {
	rd.chk = chk
}

// Reset r, bufferSize and validator of Reader
func (rd *Reader) Reset(r io.Reader, opt ...ReaderOption) {
	rd.br.Reset(r)

	if len(opt) == 0 {
		rd.dec = nil
		rd.chk = nil
	} else {
		for _, o := range opt {
			if o == nil {
				continue
			}
			o.apply(rd)
		}
	}
	if r != nil && rd.chk == nil {
		if rd.chk == nil {
			rd.chk = &DefaultEofChecker{}
		}
	}
}

// Read reads the data of underlying io.Reader into p
func (rd *Reader) Read(p []byte) (n int, err error) {
	n, err = rd.br.Read(p)
	return
}

// ReadAll reads all data from r.r and returns it.
func (rd *Reader) ReadAll() ([]byte, error) {
	return readAll(rd.br, rd.chk)
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

func (rd *Reader) Decode(v any) error {
	if rd.dec == nil {
		return ErrNoDecoder
	}
	return rd.dec.Decode(v)
}

func (rd *Reader) Peek(n int) ([]byte, error) {
	return rd.br.Peek(n)
}

func (rd *Reader) PeekRest() ([]byte, error) {
	n := rd.br.Buffered()
	return rd.Peek(n)
}

// Writer writes data to underlying Writer
type Writer struct {
	bw  *bufio.Writer
	enc Encoder
}

func newWriter(w io.Writer, opt ...WriterOption) *Writer {
	if bw, ok := w.(*bufio.Writer); ok {
		b := &Writer{bw: bw}
		for _, o := range opt {
			o.apply(b)
		}
		return b
	}
	bw := bufio.NewWriter(w)
	b := &Writer{bw: bw}
	for _, o := range opt {
		o.apply(b)
	}
	return b
}

func (wr *Writer) Reset(w io.Writer, opt ...WriterOption) {
	wr.bw.Reset(w)

	if len(opt) == 0 {
		wr.enc = nil
	} else {
		for _, o := range opt {
			if o == nil {
				continue
			}
			o.apply(wr)
		}
	}
}

func (wr *Writer) Write(p []byte) (n int, err error) {
	n, err = wr.bw.Write(p)
	return
}

func (wr *Writer) Flush() error {
	return wr.bw.Flush()
}

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

func (wr *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = wr.bw.ReadFrom(r)
	return
}

var ErrNoEncoder = errors.New("no encoder was provided")

// Encode writes encoded data into underlying Writer.
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
