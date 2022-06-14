package sicore_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/go-wonk/si/sicore"
)

func BenchmarkHmacSha256HexEncoded_Basic(b *testing.B) {
	secret := "1234"
	msg := bytes.Repeat([]byte("asdf"), 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hm := hmac.New(sha256.New, []byte(secret))
		_, _ = hm.Write(msg)
		_ = hex.EncodeToString(hm.Sum(nil))
	}
}

func BenchmarkHmacSha256HexEncoded(b *testing.B) {
	secret := "1234"
	msg := bytes.Repeat([]byte("asdf"), 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = sicore.HmacSha256HexEncoded(secret, msg)
	}
}

func BenchmarkHmacSha256HexEncoded_Basic_LargeMsg(b *testing.B) {
	secret := "1234"
	msg := bytes.Repeat([]byte("asdf"), 1000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hm := hmac.New(sha256.New, []byte(secret))
		_, _ = hm.Write(msg)
		_ = hex.EncodeToString(hm.Sum(nil))
	}
}

func BenchmarkHmacSha256HexEncoded_LargeMsg(b *testing.B) {
	secret := "1234"
	msg := bytes.Repeat([]byte("asdf"), 1000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = sicore.HmacSha256HexEncoded(secret, msg)
	}
}
