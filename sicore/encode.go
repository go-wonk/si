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

// EncodeFunc wraps a function that conforms to Encoder interface.
type EncodeFunc func(v any) error

// Encode encodes input v
func (e EncodeFunc) Encode(v any) error {
	return e(v)
}

// DefaultEncoder bypasses v to w.
// Only allows []byte or string.
func DefaultEncoder(w io.Writer) Encoder {
	return EncodeFunc(func(v any) error {
		switch c := v.(type) {
		case []byte:
			_, err := write(w, c)
			return err
		case *[]byte:
			_, err := write(w, *c)
			return err
		case string:
			_, err := write(w, []byte(c))
			return err
		case *string:
			_, err := write(w, []byte(*c))
			return err
		default:
			return errors.New("input type is not []byte")
		}
	})
}

// JsonEncoder encode v to json format and write to w
func JsonEncoder(w io.Writer) Encoder {
	return EncodeFunc(func(v any) error {
		enc := json.NewEncoder(w)
		err := enc.Encode(v)
		if err != nil {
			return err
		}
		if f, ok := w.(Flusher); ok {
			return f.Flush()
		}
		return nil
	})
}
