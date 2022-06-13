package sihttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	urlpkg "net/url"

	"github.com/go-wonk/si/sicore"
	"golang.org/x/net/http/httpguts"
)

type HttpRequest struct {
	*http.Request
	buf *bytes.Buffer
}

// newHttpRequest creates a new http.Request.
// `body` argument is not fed into NewRequest function, but is set using `setBody` after instantiating
// a http.Request.
func newHttpRequest(method string, url string, header http.Header, body any, opts ...sicore.WriterOption) (*HttpRequest, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq := HttpRequest{
		Request: req,
		buf:     bytes.NewBuffer(make([]byte, 0, 512)),
	}

	if r, ok := body.(io.Reader); ok {
		httpReq.setBody(r)
	} else {
		if err := httpReq.writeAndSetBody(body, opts...); err != nil {
			return nil, err
		}
	}

	httpReq.SetHeader(header)

	return &httpReq, nil
}

func (hr *HttpRequest) Reset(method string, url string, header http.Header, body any, opts ...sicore.WriterOption) error {

	hr.Body = nil
	hr.GetBody = nil

	if body == nil {
		hr.setBody(nil)
	} else {
		hr.buf.Reset()

		if r, ok := body.(io.Reader); ok {
			hr.setBody(r)
		} else {
			if err := hr.writeAndSetBody(body, opts...); err != nil {
				return err
			}
		}
	}

	// Set method, url, host
	if err := hr.setMethodAndURL(method, url); err != nil {
		return err
	}

	// Clear headers
	for k := range hr.Header {
		delete(hr.Header, k)
	}

	// Clear trailers
	for k := range hr.Trailer {
		delete(hr.Trailer, k)
	}

	// Clear transfer encodings
	hr.TransferEncoding = hr.TransferEncoding[:0]

	hr.SetHeader(header)

	return nil
}

// SetHeader sets `haeder` to underlying Request.
func (hr *HttpRequest) SetHeader(header http.Header) {
	for k, val := range header {
		for i, v := range val {
			if i == 0 {
				hr.Header.Set(k, v)
				continue
			}
			hr.Header.Add(k, v)
		}
	}
}

// setBody sets `body` to underlying Request.
// Most part of this function was brought from default net/http package's NewRequest function.
// It handles `sicore.Reader` and `sicore.ReadWriter`
func (hr *HttpRequest) setBody(body io.Reader) {
	hr.ContentLength = 0

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}

	hr.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			hr.ContentLength = int64(v.Len())
			buf := v.Bytes()
			hr.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(r), nil
			}
		case *bytes.Reader:
			hr.ContentLength = int64(v.Len())
			snapshot := *v
			hr.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *strings.Reader:
			hr.ContentLength = int64(v.Len())
			snapshot := *v
			hr.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *sicore.Reader:
			hr.ContentLength = int64(v.Len())
			snapshot := *v
			hr.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *sicore.ReadWriter:
			hr.ContentLength = int64(v.RLen())
			snapshot := *v
			hr.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		default:
			// This is where we'd set it to -1 (at least
			// if body != NoBody) to mean unknown, but
			// that broke people during the Go 1.8 testing
			// period. People depend on it being 0 I
			// guess. Maybe retry later. See Issue 18117.
		}
		// For client requests, Request.ContentLength of 0
		// means either actually 0, or unknown. The only way
		// to explicitly say that the ContentLength is zero is
		// to set the Body to nil. But turns out too much code
		// depends on NewRequest returning a non-nil Body,
		// so we use a well-known ReadCloser variable instead
		// and have the http package also treat that sentinel
		// variable to mean explicitly zero.
		if hr.GetBody != nil && hr.ContentLength == 0 {
			hr.Body = http.NoBody
			hr.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	} else {
		hr.GetBody = nil
		hr.Body = nil
	}
}

func (hr *HttpRequest) writeAndSetBody(body any, opts ...sicore.WriterOption) error {
	w := sicore.GetWriter(hr.buf, opts...)
	defer sicore.PutWriter(w)

	if err := w.EncodeFlush(body); err != nil {
		return err
	}

	hr.setBody(hr.buf)
	return nil
}

// setMethodAndURL sets method and url to underlying request.
func (hr *HttpRequest) setMethodAndURL(method string, url string) error {
	if !validMethod(method) {
		return fmt.Errorf("invalid method %q", method)
	}

	u, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)

	hr.Method = method
	hr.URL = u
	hr.Host = u.Host

	// hr.Proto = "HTTP/1.1"
	// hr.ProtoMajor = 1
	// hr.ProtoMinor = 1

	return nil
}

/*
Functions below are from the default package.
They are needed to create/modify the way of creating http.Request.
*/

// isNotToken is brought from default package net/http(http.go).
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

// validMethod checks whether `method` is valid.
// This function is from default package net/http(request.go).
func validMethod(method string) bool {
	/*
	     Method         = "OPTIONS"                ; Section 9.2
	                    | "GET"                    ; Section 9.3
	                    | "HEAD"                   ; Section 9.4
	                    | "POST"                   ; Section 9.5
	                    | "PUT"                    ; Section 9.6
	                    | "DELETE"                 ; Section 9.7
	                    | "TRACE"                  ; Section 9.8
	                    | "CONNECT"                ; Section 9.9
	                    | extension-method
	   extension-method = token
	     token          = 1*<any CHAR except CTLs or separators>
	*/
	return len(method) > 0 && strings.IndexFunc(method, isNotToken) == -1
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
// This function is from default package net/http(http.go).
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
// This function is from default package net/http(http.go).
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}
