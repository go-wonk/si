package sicore

import (
	"errors"
	"io"
)

const defaultBufferSize = 4096

// Flusher interface has Flush method to check if a writer has a flush method like bufio.Writer.
// json.Encoder doesn't flush after write.
type Flusher interface {
	Flush() error
}

// ReadWriter uses bufio package to read and write more efficiently.
// It is designed to read/write data from/to a storage that implements ReadWriter interface.
// `validator` determines when to finish reading, and defaultValidator is to finish when io.EOF is met.
type ReadWriter struct {
	r io.Reader
	w io.Writer
	// rw         io.ReadWriter
	bufferSize int
	validator  ReadValidator

	enc Encoder
}

// NewReadWriter with default bufferSize
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return NewReadWriterSizeWithValidatorAndEncoder(rw, defaultBufferSize, defaultValidate(), DefaultEncoder(rw))
}

// NewReadWriterWithEncoder with default bufferSize
func NewReadWriterWithEncoder(rw io.ReadWriter, enc Encoder) *ReadWriter {
	return NewReadWriterSizeWithValidatorAndEncoder(rw, defaultBufferSize, defaultValidate(), enc)
}

// NewReadWriterWithValidator with defaultBufferSize and specified validator
func NewReadWriterWithValidator(rw io.ReadWriter, validator ReadValidator) *ReadWriter {
	return NewReadWriterSizeWithValidatorAndEncoder(rw, defaultBufferSize, validator, DefaultEncoder(rw))
}

// NewReadWriterSize with specified bufferSize
func NewReadWriterSize(rw io.ReadWriter, bufferSize int) *ReadWriter {
	return NewReadWriterSizeWithValidatorAndEncoder(rw, bufferSize, defaultValidate(), DefaultEncoder(rw))
}

// NewReadWriterSizeWithValidator with specified bufferSize and validator
func NewReadWriterSizeWithValidator(rw io.ReadWriter, bufferSize int, validator ReadValidator) *ReadWriter {
	return NewReadWriterSizeWithValidatorAndEncoder(rw, bufferSize, validator, DefaultEncoder(rw))
}

// NewReadWriterSizeWithValidator with specified bufferSize and validator
func NewReadWriterSizeWithValidatorAndEncoder(rw io.ReadWriter, bufferSize int, validator ReadValidator, enc Encoder) *ReadWriter {
	return &ReadWriter{rw, rw, bufferSize, validator, enc}
}

// Read reads data into p.
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	br := getBufioReader(rw.r)
	defer putBufioReader(br)
	return br.Read(p)
}

// Write writes p into rw.w
func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	return write(rw.w, p)
}

// ReadAll reads all data in rw.r and returns it.
func (rw *ReadWriter) ReadAll() ([]byte, error) {
	return readAll(rw.r, rw.bufferSize, rw.validator)
}

// WriteAndRead writes p into rw.w then reads all data from rw.r then returns it.
func (rw *ReadWriter) WriteAndRead(p []byte) ([]byte, error) {
	if n, err := write(rw.w, p); err != nil {
		return nil, err
	} else if n != len(p) {
		return nil, errors.New("bytes to write differ from what has been written")
	}

	return readAll(rw.r, rw.bufferSize, rw.validator)
}

func (rw *ReadWriter) WriteAny(p any) (n int, err error) {
	if err = rw.enc.Encode(p); err != nil {
		n = 0
		return
	}
	n = 1
	return
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

// Reader
type Reader struct {
	r          io.Reader
	bufferSize int
	validator  ReadValidator
}

func newReader() *Reader {
	return &Reader{}
}

// Reset r, bufferSize and validator of Reader
func (rd *Reader) Reset(r io.Reader, bufferSize int, validator ReadValidator) {
	rd.r = r
	rd.bufferSize = bufferSize
	rd.validator = validator
}

// Read reads the data of underlying io.Reader into p
func (r *Reader) Read(p []byte) (n int, err error) {
	br := getBufioReader(r.r)
	defer putBufioReader(br)
	return br.Read(p)
}

// ReadAll reads all data from r.r and returns it.
func (r *Reader) ReadAll() ([]byte, error) {
	return readAll(r.r, r.bufferSize, r.validator)
}
