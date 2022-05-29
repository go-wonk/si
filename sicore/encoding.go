package sicore

import (
	"encoding/json"
	"errors"
	"io"
)

// Encoder encode input parameter and write to a writer.
type Encoder interface {
	Encode(v any) error
}

func SetJsonEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.enc = json.NewEncoder(w)
	})
}

type DefaultEncoder struct {
	w io.Writer
}

func (de *DefaultEncoder) Encode(v any) error {
	switch c := v.(type) {
	case []byte:
		_, err := de.w.Write(c)
		return err
	case *[]byte:
		_, err := de.w.Write(*c)
		return err
	case string:
		_, err := de.w.Write([]byte(c))
		return err
	case *string:
		_, err := de.w.Write([]byte(*c))
		return err
	default:
		return errors.New("input type is not allowed")
	}

}

func SetDefaultEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.enc = &DefaultEncoder{w}
	})
}

type Decoder interface {
	Decode(v any) error
}

func SetJsonDecoder() ReaderOption {
	return ReaderOptionFunc(func(r *Reader) {
		r.dec = json.NewDecoder(r)
	})
}
