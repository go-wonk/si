package sihttp

import (
	"net/http"
	"strings"

	"github.com/go-wonk/si/sicore"
)

type RequestOption interface {
	apply(c *http.Request) error
}

type RequestOptionFunc func(c *http.Request) error

func (o RequestOptionFunc) apply(c *http.Request) error {
	return o(c)
}

func WithHeaderHmac256(key string, secret []byte) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		header := req.Header
		if _, ok := header[key]; ok {
			// skip
			return nil
		}

		contentType := header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// skip
			return nil
		}

		if req.GetBody == nil {
			// skip
			return nil
		}

		r, err := req.GetBody()
		if err != nil {
			return err
		}

		hashed, err := sicore.HmacSha256HexEncodedWithReader(string(secret), r)
		if err != nil {
			return err
		}
		header[key] = []string{hashed}

		return nil
	})
}

// type RequestOption interface {
// 	apply(c *HttpClient)
// }
// type RequestOptionFunc func(c *HttpClient)

// func (o RequestOptionFunc) apply(c *HttpClient) {
// 	o(c)
// }

// // WithEncoder sets Client's encoder
// func WithEncoder(enc sicore.Encoder) RequestOptionFunc {
// 	return RequestOptionFunc(func(c *HttpClient) {
// 		// c.SetEncoder(enc)
// 	})
// }

// func WithJsonEncoder(w io.Writer) RequestOptionFunc {
// 	return RequestOptionFunc(func(c *HttpClient) {
// 		// enc := json.NewEncoder(w)
// 		// c.SetEncoder(enc)
// 	})
// }
