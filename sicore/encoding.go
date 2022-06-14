package sicore

import (
	"encoding/hex"
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
	if v == nil {
		return nil
	}

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

// HmacSha256HexEncoded returns hmac sha256 hash's hex string.
func HmacSha256HexEncoded(secret string, message []byte) (string, error) {
	hm := GetHmacSha256Hash(secret)
	defer PutHmacSha256Hash(secret, hm)
	_, err := hm.Write(message)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hm.Sum(nil)), nil
}

// HmacSha256HexEncodedWithReader returns hmac sha256 hash's hex string.
func HmacSha256HexEncodedWithReader(secret string, r io.Reader) (string, error) {
	body, err := ReadAll(r)
	if err != nil {
		return "", err
	}
	return HmacSha256HexEncoded(secret, body)
}
