package sihttp

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/go-wonk/si/v2/sicore"
	"golang.org/x/oauth2"
)

type RequestOption interface {
	apply(c *http.Request) error
}

type RequestOptionFunc func(c *http.Request) error

func (o RequestOptionFunc) apply(c *http.Request) error {
	return o(c)
}

func WithHeaderSet(key string, value string) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		header := req.Header
		header[key] = []string{value}
		return nil
	})
}

func WithHeaderAdd(key string, value string) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		header := req.Header
		header.Add(key, value)
		return nil
	})
}

func WithBasicAuth(username, password string) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		header := req.Header
		header["Authorization"] = []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}
		return nil
	})
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

func WithTokenSource(tokenSource oauth2.TokenSource) RequestOptionFunc {
	return RequestOptionFunc(func(req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return err
		}
		header := req.Header
		if _, ok := header["Authorization"]; ok {
			// skip
			return nil
		}
		header["Authorization"] = []string{token.TokenType + " " + token.AccessToken}
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
	apply(c *Client) error
}

type ClientOptionFunc func(c *Client) error

func (o ClientOptionFunc) apply(c *Client) error {
	return o(c)
}

func WithBaseUrl(baseUrl string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.baseUrl = baseUrl
		return nil
	})
}

// WithDefaultHeaders set Client with defaultHeaders that will be set on every request
func WithDefaultHeaders(defaultHeaders map[string]string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.defaultHeaders = defaultHeaders
		return nil
	})
}

func WithReaderOpt(opt sicore.ReaderOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.appendReaderOption(opt)
		return nil
	})
}

func WithWriterOpt(opt sicore.WriterOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.appendWriterOption(opt)
		return nil
	})
}

func WithRequestOpt(opt RequestOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.appendRequestOption(opt)
		return nil
	})
}

func WithRequestHeaderHmac256(key string, secret []byte) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.appendRequestOption(WithHeaderHmac256(key, secret))
		return nil
	})
}

func WithRetryAttempts(attempts int) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) error {
		c.retryAttempts = attempts
		return nil
	})
}
