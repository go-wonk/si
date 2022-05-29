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
	br        *bufio.Reader
	validator ReadValidator
}

func newReader(r io.Reader, val ReadValidator) *Reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &Reader{br, val}
	}
	br := bufio.NewReader(r)
	return &Reader{br, val}
}

// Reset r, bufferSize and validator of Reader
func (rd *Reader) Reset(r io.Reader, validator ReadValidator) {
	rd.br.Reset(r)
	rd.validator = validator
}

// Read reads the data of underlying io.Reader into p
func (rd *Reader) Read(p []byte) (n int, err error) {
	n, err = rd.br.Read(p)
	return
}

// ReadAll reads all data from r.r and returns it.
func (rd *Reader) ReadAll() ([]byte, error) {
	return readAll(rd.br, rd.validator)
}

// readAll reads all data from r and returns it
func readAll(r io.Reader, validator ReadValidator) ([]byte, error) {

	b := make([]byte, 0, defaultBufferSize)
	for {
		if len(b) == cap(b) {
			if err := growCap(&b, defaultBufferSize); err != nil {
				return nil, err
			}
		}

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		ended, err := validator.validate(b, err)
		if err != nil {
			return b, err
		}
		if ended {
			return b, nil
		}
	}
}

// Writer writes data to underlying Writer
type Writer struct {
	bw  *bufio.Writer
	enc Encoder
}

func newWriter(w io.Writer, opt ...Option) *Writer {
	if bw, ok := w.(*bufio.Writer); ok {
		b := &Writer{bw: bw}
		for _, o := range opt {
			o.Apply(b)
		}
		return b
	}
	bw := bufio.NewWriter(w)
	b := &Writer{bw: bw}
	for _, o := range opt {
		o.Apply(b)
	}
	return b
}

func (wr *Writer) Reset(w io.Writer, opt ...Option) {
	wr.bw.Reset(w)

	if len(opt) == 0 {
		wr.enc = nil
	} else {
		if w != nil {
			for _, o := range opt {
				o.Apply(wr)
			}
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
// It flushes any remaining buffer right away.
func (wr *Writer) Encode(p any) (err error) {
	if wr.enc == nil {
		return ErrNoEncoder
	}
	if err = wr.enc.Encode(p); err != nil {
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

func NewReadWriter(r *Reader, w *Writer) *ReadWriter {
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
	if err := rw.Encode(v); err != nil {
		return nil, err
	}
	b, err := rw.ReadAll()
	if err != nil {
		return nil, err
	}
	return b, nil
}
