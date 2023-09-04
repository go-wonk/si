package siutils_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/siutils"
	"github.com/go-wonk/si/v2/tests/testmodels"
	"github.com/stretchr/testify/assert"
)

func TestHMacSha256(t *testing.T) {

	secret := []byte("1234")
	mac := hmac.New(sha256.New, secret)

	message := []byte("my message")
	_, err := mac.Write(message)
	siutils.AssertNilFail(t, err)

	hashed := hex.EncodeToString(mac.Sum(nil))
	// fmt.Println(hashed)
	assert.EqualValues(t, "34420f26f2612cb4e0a812c5e39f656390e4f6c91699d44303425e37bb979d0a", hashed)
}

func TestHMacSha256_JsonString(t *testing.T) {

	secret := []byte("1234")
	mac := hmac.New(sha256.New, secret)

	message := []byte(`{"id":1,"email_address":"wonk@wonk.org","name":"wonk","borrowed":true}` + "\n")
	_, err := mac.Write(message)
	siutils.AssertNilFail(t, err)

	hashed := hex.EncodeToString(mac.Sum(nil))
	// fmt.Println(hashed)
	// assert.EqualValues(t, "21f8ffcfa671fce6dda76323af0b16c07f7d95242b30fb34e3f0e2a202cd8784", hashed)
	assert.EqualValues(t, "8af0704ab60a967daa7bef4e6e9d2add957c0cec36b5c98b9810e2a6d8ebae30", hashed)
}

func TestHMacSha256_Writer(t *testing.T) {

	secret := []byte("1234")
	message := []byte("my message")

	mac := hmac.New(sha256.New, secret)
	w := sicore.GetWriter(mac)
	defer sicore.PutWriter(w)

	_, err := w.WriteFlush(message)
	siutils.AssertNilFail(t, err)

	hashed := hex.EncodeToString(mac.Sum(nil))
	assert.EqualValues(t, "34420f26f2612cb4e0a812c5e39f656390e4f6c91699d44303425e37bb979d0a", hashed)
}

func TestHMacSha256_Writer_Encode(t *testing.T) {

	secret := []byte("1234")
	s := testmodels.Student{ID: 1, EmailAddress: "wonk@wonk.org", Name: "wonk", Borrowed: true}

	mac := hmac.New(sha256.New, secret)
	w := sicore.GetWriter(mac, sicore.SetJsonEncoder())
	defer sicore.PutWriter(w)

	err := w.EncodeFlush(&s)
	siutils.AssertNilFail(t, err)

	hashed := hex.EncodeToString(mac.Sum(nil))
	assert.EqualValues(t, "8af0704ab60a967daa7bef4e6e9d2add957c0cec36b5c98b9810e2a6d8ebae30", hashed)
}

func BenchmarkHMacSha256(b *testing.B) {
	secret := []byte("1234")
	message := []byte("my message")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mac := hmac.New(sha256.New, secret)
		_, err := mac.Write(message)
		siutils.AssertNilFailB(b, err)

		_ = hex.EncodeToString(mac.Sum(nil))
	}

}
func BenchmarkHMacSha256_Writer(b *testing.B) {
	secret := []byte("1234")
	message := []byte("my message")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mac := hmac.New(sha256.New, secret)
		w := sicore.GetWriter(mac)

		_, err := w.WriteFlush(message)
		siutils.AssertNilFailB(b, err)

		_ = hex.EncodeToString(mac.Sum(nil))
		sicore.PutWriter(w)
	}

}

func BenchmarkHMacSha256_JsonEncode(b *testing.B) {
	secret := []byte("1234")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mac := hmac.New(sha256.New, secret)

		s := testmodels.Student{ID: 1, EmailAddress: "wonk@wonk.org", Name: "wonk", Borrowed: true}

		_, err := mac.Write([]byte(s.String() + "\n"))
		siutils.AssertNilFailB(b, err)

		_ = hex.EncodeToString(mac.Sum(nil))
	}

}

func BenchmarkHMacSha256_Writer_JsonEncode(b *testing.B) {
	secret := []byte("1234")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mac := hmac.New(sha256.New, secret)

		s := testmodels.Student{ID: 1, EmailAddress: "wonk@wonk.org", Name: "wonk", Borrowed: true}
		w := sicore.GetWriter(mac, sicore.SetJsonEncoder())

		err := w.EncodeFlush(&s)
		siutils.AssertNilFailB(b, err)

		_ = hex.EncodeToString(mac.Sum(nil))
		sicore.PutWriter(w)
	}

}

func BenchmarkHMacSha256_Writer_JsonMarshal(b *testing.B) {
	secret := []byte("1234")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mac := hmac.New(sha256.New, secret)

		s := testmodels.Student{ID: 1, EmailAddress: "wonk@wonk.org", Name: "wonk", Borrowed: true}
		j := s.String() + "\n"
		w := sicore.GetWriter(mac)

		_, err := w.WriteFlush([]byte(j))
		siutils.AssertNilFailB(b, err)

		_ = hex.EncodeToString(mac.Sum(nil))
		sicore.PutWriter(w)
	}

}
