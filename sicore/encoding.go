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

type Option interface {
	Apply(w *Writer)
}

type OptionFunc func(*Writer)

func (s OptionFunc) Apply(w *Writer) {
	s(w)
}

func SetJsonEncoder() Option {
	return OptionFunc(func(w *Writer) {
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

func SetDefaultEncoder() Option {
	return OptionFunc(func(w *Writer) {
		w.enc = &DefaultEncoder{w}
	})
}
