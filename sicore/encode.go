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

// // EncodeFunc wraps a function that conforms to Encoder interface.
// type EncodeFunc func(v any) error

// // Encode encodes input v
// func (e EncodeFunc) Encode(v any) error {
// 	return e(v)
// }

// // DefaultEncoder bypasses v to w.
// // Only allows []byte or string.
// func DefaultEncoder(w io.Writer) Encoder {
// 	return EncodeFunc(func(v any) error {
// 		switch c := v.(type) {
// 		case []byte:
// 			_, err := write(w, c)
// 			return err
// 		case *[]byte:
// 			_, err := write(w, *c)
// 			return err
// 		case string:
// 			_, err := write(w, []byte(c))
// 			return err
// 		case *string:
// 			_, err := write(w, []byte(*c))
// 			return err
// 		default:
// 			return errors.New("input type is not []byte")
// 		}
// 	})
// }

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
