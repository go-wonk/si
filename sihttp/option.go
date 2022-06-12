package sihttp

import (
	"io"

	"github.com/go-wonk/si/sicore"
)

type RequestOption interface {
	apply(c *HttpClient)
}
type RequestOptionFunc func(c *HttpClient)

func (o RequestOptionFunc) apply(c *HttpClient) {
	o(c)
}

// WithEncoder sets Client's encoder
func WithEncoder(enc sicore.Encoder) RequestOptionFunc {
	return RequestOptionFunc(func(c *HttpClient) {
		// c.SetEncoder(enc)
	})
}

func WithJsonEncoder(w io.Writer) RequestOptionFunc {
	return RequestOptionFunc(func(c *HttpClient) {
		// enc := json.NewEncoder(w)
		// c.SetEncoder(enc)
	})
}
