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

func WithBearerToken(token string) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		header := req.Header
		if _, ok := header["Authorization"]; ok {
			// skip
			return nil
		}
		header["Authorization"] = []string{"Bearer " + token}
		return nil
	})
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

type ClientOption interface {
	apply(c *HttpClient) error
}

type ClientOptionFunc func(c *HttpClient) error

func (o ClientOptionFunc) apply(c *HttpClient) error {
	return o(c)
}

func WithReaderOpt(opt sicore.ReaderOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *HttpClient) error {
		c.appendReaderOption(opt)
		return nil
	})
}

func WithWriterOpt(opt sicore.WriterOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *HttpClient) error {
		c.appendWriterOption(opt)
		return nil
	})
}

func WithRequestOpt(opt RequestOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *HttpClient) error {
		c.appendRequestOption(opt)
		return nil
	})
}

func WithRequestHeaderHmac256(key string, secret []byte) ClientOptionFunc {
	return ClientOptionFunc(func(c *HttpClient) error {
		c.appendRequestOption(WithHeaderHmac256(key, secret))
		return nil
	})
}
