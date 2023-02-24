package sihttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// DefaultInsecureStandardClient instantiate http.Client with InsecureSkipVerify set to true
func DefaultInsecureStandardClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return DefaultStandardClient(tlsConfig)
}

// DefaultStandardClient instantiate http.Client with input parameter `tlsConfig`
func DefaultStandardClient(tlsConfig *tls.Config) *http.Client {

	dialer := &net.Dialer{Timeout: 5 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       50,
		IdleConnTimeout:    time.Duration(60) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    tlsConfig,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	return NewStandardClient(30*time.Second, tr)
}

func NewStandardClient(clientTimeout time.Duration, transport *http.Transport) *http.Client {

	client := &http.Client{
		Timeout:   time.Duration(30) * time.Second,
		Transport: transport,
	}
	return client
}
