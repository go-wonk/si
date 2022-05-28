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

type EncoderSetter interface {
	SetEncoder(w *Writer)
}

type SetEncoderFunc func(*Writer)

func (s SetEncoderFunc) SetEncoder(w *Writer) {
	s(w)
}

func SetJsonEncoder() EncoderSetter {
	return SetEncoderFunc(func(w *Writer) {
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

func SetDefaultEncoder() EncoderSetter {
	return SetEncoderFunc(func(w *Writer) {
		w.enc = &DefaultEncoder{w}
	})
}
