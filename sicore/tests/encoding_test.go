package sicore_test

import (
	"bytes"
	"testing"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/siutils"
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

func TestEncoding_BytesDecoder_DecodeBytes(t *testing.T) {
	r := bytes.NewReader([]byte("hey"))
	d := sicore.NewDefaultDecoder(r)

	var b []byte
	err := d.Decode(&b)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, []byte("hey"), b)
}

func TestEncoding_BytesDecoder_DecodeString(t *testing.T) {
	r := bytes.NewReader([]byte("hey"))
	d := sicore.NewDefaultDecoder(r)

	var s string
	err := d.Decode(&s)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, []byte("hey"), []byte(s))
}
