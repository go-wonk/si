package sicore

import (
	"errors"
	"io"
)

const defaultBufferSize = 4096

// File, Tcp, Ftp
type ReadWriter struct {
	r io.Reader
	w io.Writer
	// rw         io.ReadWriter
	bufferSize int
	validator  ReadValidator
}

func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return NewReadWriterSize(rw, defaultBufferSize)
}

func NewReadWriterWithValidator(rw io.ReadWriter, validator ReadValidator) *ReadWriter {
	return NewReadWriterSizeWithValidator(rw, defaultBufferSize, validator)
}

func NewReadWriterSize(rw io.ReadWriter, bufferSize int) *ReadWriter {
	return NewReadWriterSizeWithValidator(rw, bufferSize, defaultValidate())
}

func NewReadWriterSizeWithValidator(rw io.ReadWriter, bufferSize int, validator ReadValidator) *ReadWriter {
	return &ReadWriter{rw, rw, bufferSize, validator}
}

// func (rw *ReadWriter) SetValidator(validator ReadValidator) {
// 	rw.validator = validator
// }

func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	br := getBufioReader(rw.r)
	defer putBufioReader(br)
	return br.Read(p)
}

func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	return rw.write(p)
}

func (rw *ReadWriter) ReadAll() ([]byte, error) {
	return readAll(rw.r, rw.bufferSize, rw.validator)
}

func (rw *ReadWriter) WriteAndRead(p []byte) ([]byte, error) {
	if n, err := rw.write(p); err != nil {
		return nil, err
	} else if n != len(p) {
		return nil, errors.New("bytes to write differ from what has been written")
	}

	return readAll(rw.r, rw.bufferSize, rw.validator)
}

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
