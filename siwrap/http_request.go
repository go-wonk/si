package siwrap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	urlpkg "net/url"

	"github.com/go-wonk/si/sicore"
	"golang.org/x/net/http/httpguts"
)

//
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}
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
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func setHeader(req *http.Request, header http.Header) {
	for k, val := range header {
		for i, v := range val {
			if i == 0 {
				req.Header.Set(k, v)
				continue
			}
			req.Header.Add(k, v)
		}
	}
}

func setBody(req *http.Request, body io.Reader) {
	req.ContentLength = 0

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}

	req.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *sicore.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *sicore.ReadWriter:
			req.ContentLength = int64(v.RLen())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
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
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
}

func setMethodAndURL(req *http.Request, method string, url string) error {
	if !validMethod(method) {
		return fmt.Errorf("invalid method %q", method)
	}

	u, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)

	req.Method = method
	req.URL = u
	req.Host = u.Host

	// req.Proto = "HTTP/1.1"
	// req.ProtoMajor = 1
	// req.ProtoMinor = 1

	return nil
}

func newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set Body, GetBody(), Content Length
	setBody(req, body)

	return req, nil
}

// pool
var (
	_requestPool = sync.Pool{}
)

// GetRequest retrieves a request from a pool and returns it.
func GetRequest(method string, url string, body io.Reader) (*http.Request, error) {
	g := _requestPool.Get()
	if g == nil {
		return newRequest(method, url, body)
	}
	req := g.(*http.Request)

	// Set Body, GetBody(), Content Length
	setBody(req, body)

	// Set method, url, host
	if err := setMethodAndURL(req, method, url); err != nil {
		return nil, err
	}

	// Clear headers
	for k := range req.Header {
		delete(req.Header, k)
	}

	// Clear trailers
	for k := range req.Trailer {
		delete(req.Trailer, k)
	}

	// Clear transfer encodings
	req.TransferEncoding = req.TransferEncoding[:0]

	return req, nil
}

func PutRequest(req *http.Request) {
	req.Body.Close()
	_requestPool.Put(req)
}
