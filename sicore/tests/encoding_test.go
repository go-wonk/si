package sicore_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/go-wonk/si"
	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

func TestHmacSha256HexEncoded(t *testing.T) {
	secret := "1234"
	hashed, err := sicore.HmacSha256HexEncoded(secret, []byte("asdf"))
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "e5e9f44b2dcbe23988aa01743748a5fe64f494d7c5eea29ea94ae4e34878868e", hashed)

	hashed, err = sicore.HmacSha256HexEncoded(secret, []byte("qwer"))
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "685f4fdb529e85b9e8fab7f9daaf550b5534e956d5c5f0f7a33c1ade0d8d67ea", hashed)
}

type NopDecoder struct {
	r io.Reader
}

func (d *NopDecoder) Decode(v any) error {
	switch t := v.(type) {
	case *[]byte:
		b, err := si.ReadAll(d.r)
		if err != nil {
			return err
		}
		*t = b
		return nil
	}

	return errors.New("not supported")
}

func TestDecodeAny(t *testing.T) {
	r := bytes.NewReader([]byte("hey"))
	d := NopDecoder{r}
	var b []byte
	err := d.Decode(&b)
	siutils.AssertNilFail(t, err)

	// log.Println(string(b))
	assert.EqualValues(t, []byte("hey"), b)
}

func BenchmarkDecodeAny_PassByPointer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader([]byte("hey"))
		d := NopDecoder{r}
		var byt []byte
		err := d.Decode(&byt)
		siutils.AssertNilFailB(b, err)
	}
}

func BenchmarkDecodeAny_Return(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader([]byte("hey"))
		_, err := si.ReadAll(r)
		siutils.AssertNilFailB(b, err)
	}
}
