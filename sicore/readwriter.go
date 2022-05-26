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

func (rw *ReadWriter) ReadAllBytes() ([]byte, error) {
	return rw.readAll()

}

func (rw *ReadWriter) WriteAndRead(p []byte) ([]byte, error) {
	if n, err := rw.write(p); err != nil {
		return nil, err
	} else if n != len(p) {
		return nil, errors.New("bytes to write differ from what has been written")
	}

	return rw.readAll()
}

func (rw *ReadWriter) readAll() ([]byte, error) {
	br := getBufioReader(rw.r)
	defer putBufioReader(br)

	b := make([]byte, 0, rw.bufferSize)
	for {
		if len(b) == cap(b) {
			if err := growCap(&b, rw.bufferSize); err != nil {
				return nil, err
			}
		}

		n, err := br.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		ended, err := rw.validator.validate(b, err)
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
