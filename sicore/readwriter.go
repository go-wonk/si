package sicore

import (
	"errors"
	"io"
)

const defaultBufferSize = 4096

// ReadWriter uses bufio package to read and write more efficiently.
// It is designed to read/write data from/to a storage that implements ReadWriter interface.
// `validator` determines when to finish reading, and defaultValidator is to finish when io.EOF is met.
type ReadWriter struct {
	r io.Reader
	w io.Writer
	// rw         io.ReadWriter
	bufferSize int
	validator  ReadValidator
}

// NewReadWriter with default bufferSize
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return NewReadWriterSize(rw, defaultBufferSize)
}

// NewReadWriterWithValidator with defaultBufferSize and specified validator
func NewReadWriterWithValidator(rw io.ReadWriter, validator ReadValidator) *ReadWriter {
	return NewReadWriterSizeWithValidator(rw, defaultBufferSize, validator)
}

// NewReadWriterSize with specified bufferSize
func NewReadWriterSize(rw io.ReadWriter, bufferSize int) *ReadWriter {
	return NewReadWriterSizeWithValidator(rw, bufferSize, defaultValidate())
}

// NewReadWriterSizeWithValidator with specified bufferSize and validator
func NewReadWriterSizeWithValidator(rw io.ReadWriter, bufferSize int, validator ReadValidator) *ReadWriter {
	return &ReadWriter{rw, rw, bufferSize, validator}
}

// Read reads data into p.
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	br := getBufioReader(rw.r)
	defer putBufioReader(br)
	return br.Read(p)
}

// Write writes p into rw.w
func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	return rw.write(p)
}

// ReadAll reads all data in rw.r and returns it.
func (rw *ReadWriter) ReadAll() ([]byte, error) {
	return readAll(rw.r, rw.bufferSize, rw.validator)
}

// WriteAndRead writes p into rw.w then reads all data from rw.r then returns it.
func (rw *ReadWriter) WriteAndRead(p []byte) ([]byte, error) {
	if n, err := rw.write(p); err != nil {
		return nil, err
	} else if n != len(p) {
		return nil, errors.New("bytes to write differ from what has been written")
	}

	return readAll(rw.r, rw.bufferSize, rw.validator)
}

// readAll reads all data from r and returns it
func readAll(r io.Reader, bufferSize int, validator ReadValidator) ([]byte, error) {
	br := getBufioReader(r)
	defer putBufioReader(br)

	b := make([]byte, 0, bufferSize)
	for {
		if len(b) == cap(b) {
			if err := growCap(&b, bufferSize); err != nil {
				return nil, err
			}
		}

		n, err := br.Read(b[len(b):cap(b)])
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

func (rw *ReadWriter) write(p []byte) (int, error) {
	bw := getBufioWriter(rw.w)
	defer putBufioWriter(bw)
	n, err := bw.Write(p)
	if err != nil {
		return 0, err
	}

	if err = bw.Flush(); err != nil {
		return 0, err
	}

	return n, nil
}

// Reader
type Reader struct {
	r          io.Reader
	bufferSize int
	validator  ReadValidator
}

func newReader() *Reader {
	return &Reader{}
}

func (rd *Reader) Reset(r io.Reader, bufferSize int, validator ReadValidator) {
	rd.r = r
	rd.bufferSize = bufferSize
	rd.validator = validator
}
func (r *Reader) Read(p []byte) (n int, err error) {
	br := getBufioReader(r.r)
	defer putBufioReader(br)
	return br.Read(p)
}

func (r *Reader) ReadAll() ([]byte, error) {
	return readAll(r.r, r.bufferSize, r.validator)
}
