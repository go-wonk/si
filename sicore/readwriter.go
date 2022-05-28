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

// // ReadWriter uses bufio package to read and write more efficiently.
// // It is designed to read/write data from/to a storage that implements ReadWriter interface.
// // `validator` determines when to finish reading, and defaultValidator is to finish when io.EOF is met.
// type ReadWriter struct {
// 	r          io.Reader
// 	w          io.Writer
// 	bufferSize int
// 	validator  ReadValidator
// 	enc        Encoder
// }

// // NewReadWriter with default bufferSize
// func NewReadWriter(r io.Reader, w io.Writer) *ReadWriter {
// 	return NewReadWriterSizeWithValidatorAndEncoder(r, w, defaultBufferSize, defaultValidate(), DefaultEncoder(w))
// }

// // NewReadWriterWithEncoder with default bufferSize
// func NewReadWriterWithEncoder(r io.Reader, w io.Writer, enc Encoder) *ReadWriter {
// 	return NewReadWriterSizeWithValidatorAndEncoder(r, w, defaultBufferSize, defaultValidate(), enc)
// }

// // NewReadWriterWithValidator with defaultBufferSize and specified validator
// func NewReadWriterWithValidator(r io.Reader, w io.Writer, validator ReadValidator) *ReadWriter {
// 	return NewReadWriterSizeWithValidatorAndEncoder(r, w, defaultBufferSize, validator, DefaultEncoder(w))
// }

// // NewReadWriterSize with specified bufferSize
// func NewReadWriterSize(r io.Reader, w io.Writer, bufferSize int) *ReadWriter {
// 	return NewReadWriterSizeWithValidatorAndEncoder(r, w, bufferSize, defaultValidate(), DefaultEncoder(w))
// }

// // NewReadWriterSizeWithValidator with specified bufferSize and validator
// func NewReadWriterSizeWithValidator(r io.Reader, w io.Writer, bufferSize int, validator ReadValidator) *ReadWriter {
// 	return NewReadWriterSizeWithValidatorAndEncoder(r, w, bufferSize, validator, DefaultEncoder(w))
// }

// // NewReadWriterSizeWithValidatorAndEncoder with specified bufferSize and validator
// func NewReadWriterSizeWithValidatorAndEncoder(r io.Reader, w io.Writer, bufferSize int, validator ReadValidator, enc Encoder) *ReadWriter {
// 	return &ReadWriter{r, w, bufferSize, validator, enc}
// }

// // Read reads data into p.
// func (rw *ReadWriter) Read(p []byte) (n int, err error) {
// 	br := getBufioReader(rw.r)
// 	defer putBufioReader(br)
// 	return br.Read(p)
// }

// // Write writes p into rw.w
// func (rw *ReadWriter) Write(p []byte) (n int, err error) {
// 	return write(rw.w, p)
// }

// // ReadAll reads all data in rw.r and returns it.
// func (rw *ReadWriter) ReadAll() ([]byte, error) {
// 	return readAll(rw.r, rw.bufferSize, rw.validator)
// }

// // WriteAndRead writes p into rw.w then reads all data from rw.r then returns it.
// func (rw *ReadWriter) WriteAndRead(p []byte) ([]byte, error) {
// 	if n, err := write(rw.w, p); err != nil {
// 		return nil, err
// 	} else if n != len(p) {
// 		return nil, errors.New("bytes to write differ from what has been written")
// 	}

// 	return readAll(rw.r, rw.bufferSize, rw.validator)
// }

// func (rw *ReadWriter) WriteAny(p any) (n int, err error) {
// 	if err = rw.enc.Encode(p); err != nil {
// 		n = 0
// 		return
// 	}
// 	if f, ok := rw.w.(Flusher); ok {
// 		if err = f.Flush(); err != nil {
// 			n = 0
// 			return
// 		}
// 	}
// 	n = 1
// 	return
// }

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

// Writer writes data to underlying Writer
type Writer struct {
	bw  *bufio.Writer
	enc Encoder
}

func newWriter(w io.Writer, enc EncoderSetter) *Writer {
	if bw, ok := w.(*bufio.Writer); ok {
		b := &Writer{bw: bw}
		if enc != nil {
			enc.SetEncoder(b)
		}
		return b
	}
	bw := bufio.NewWriter(w)
	b := &Writer{bw: bw}
	if enc != nil {
		enc.SetEncoder(b)
	}
	return b
}

func (wr *Writer) Reset(w io.Writer, enc EncoderSetter) {
	wr.bw.Reset(w)
	if enc != nil {
		enc.SetEncoder(wr)
	} else {
		wr.enc = nil
	}
}

func (wr *Writer) Write(p []byte) (n int, err error) {
	n, err = wr.bw.Write(p)
	return
}

func (wr *Writer) Flush() error {
	return wr.bw.Flush()
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
	if err = wr.enc.Encode(p); err != nil {
		return
	}
	err = wr.Flush()
	return
}
