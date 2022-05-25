package sicore

import (
	"errors"
	"io"
)

const defaultBufferSize = 4096

const maxInt = int(^uint(0) >> 1)

var ErrTooLarge = errors.New("buf too large")

func grow(b *[]byte, n int) (int, error) {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		*b = (*b)[:l+n]
		return l, nil
	}

	if l+n <= c {
		// if needed length is lte c
		return l, nil
	}

	if c > maxInt-c-n {
		// too large
		return l, ErrTooLarge
	}

	newBuf := make([]byte, c*2+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l+n]
	return l, nil
}

func growCap(b *[]byte, n int) error {
	c := cap(*b)
	l := len(*b)
	a := c - l // available
	if n <= a {
		return nil
	}

	if l+n <= c {
		// if needed length is lte c
		return nil
	}

	if c > maxInt-c-n {
		// too large
		return ErrTooLarge
	}

	newBuf := make([]byte, c+n)
	copy(newBuf, (*b)[0:])
	*b = newBuf[:l]
	return nil
}

func readAll(r io.Reader, bufferSize int, validate validateFunc) ([]byte, error) {
	br := getBufioReader(r)
	defer putBufioReader(br)

	b := make([]byte, 0, bufferSize)
	for {
		if len(b) == cap(b) {
			if err := growCap(&b, bufferSize); err != nil {
				return nil, err
			}
		}

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		ended, err := validate(b, err)
		if err != nil {
			return b, err
		}
		if ended {
			return b, nil
		}
	}
}

func write(w io.Writer, p []byte) (int, error) {
	bw := getBufioWriter(w)
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

// File, Tcp, Ftp
type BytesReadWriter struct {
	rw         io.ReadWriter
	bufferSize int
	validate   validateFunc
}

func NewBytesReadWriter(rw io.ReadWriter) *BytesReadWriter {
	return NewBytesReadWriterSize(rw, defaultBufferSize)
}

func NewBytesReadWriterSize(rw io.ReadWriter, bufferSize int) *BytesReadWriter {
	return &BytesReadWriter{rw, bufferSize, defaultValidate}
}

func (rw *BytesReadWriter) SetValidateFunc(validate validateFunc) {
	rw.validate = validate
}

func (rw *BytesReadWriter) Read(p []byte) (n int, err error) {
	return rw.rw.Read(p)
}

func (rw *BytesReadWriter) ReadAllBytes() ([]byte, error) {
	return readAll(rw.rw, rw.bufferSize, rw.validate)

}

func (rw *BytesReadWriter) Write(p []byte) (n int, err error) {
	return write(rw.rw, p)
}

func (rw *BytesReadWriter) WriteAndRead(p []byte) ([]byte, error) {
	if n, err := write(rw.rw, p); err != nil {
		return nil, err
	} else if n != len(p) {
		return nil, errors.New("bytes to write differ from what has been written")
	}

	return readAll(rw.rw, rw.bufferSize, rw.validate)
}
