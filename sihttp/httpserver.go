package sihttp

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	TLSConf *tls.Config
	Server  *http.Server
	pem     string
	key     string
}

func NewHttpServer(tlsConfig *tls.Config,
	addr string, writeTimeOutInSec int, readTimeOutInSec int,
	pem string, key string) *HttpServer {

	return NewHttpServerCors(tlsConfig, addr, writeTimeOutInSec, readTimeOutInSec, pem, key, nil, nil, nil)

}

func NewHttpServerCors(tlsConfig *tls.Config,
	addr string, writeTimeOutInSec int, readTimeOutInSec int,
	pem string, key string,
	allowedOrigins, allowedHeaders, allowedMethods []string) *HttpServer {

	var hs HttpServer
	hs.TLSConf = tlsConfig
	hs.Server = &http.Server{
		Addr:         addr,
		TLSConfig:    hs.TLSConf,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		WriteTimeout: time.Duration(writeTimeOutInSec) * time.Second,
		ReadTimeout:  time.Duration(readTimeOutInSec) * time.Second,
	}
	hs.pem = pem
	hs.key = key

	router := mux.NewRouter()
	if len(allowedOrigins) == 0 && len(allowedHeaders) == 0 && len(allowedMethods) == 0 {
		hs.Server.Handler = router
	} else {
		corsOrigin := handlers.AllowedOrigins(allowedOrigins)
		corsHeaders := handlers.AllowedHeaders(allowedHeaders)
		corsMethods := handlers.AllowedMethods(allowedMethods)
		// nr := handlers.CORS(corsOrigin, corsHeaders, corsMethods)(hs.Router)
		cors := handlers.CORS(corsOrigin, corsHeaders, corsMethods)
		corsHandler := cors(router)

		hs.Server.Handler = corsHandler
	}

	return &hs
}

func (hs *HttpServer) Start() error {
	var err error
	if len(hs.pem) > 0 || len(hs.key) > 0 {
		err = hs.Server.ListenAndServeTLS(hs.pem, hs.key)
	} else {
		err = hs.Server.ListenAndServe()
	}
	return err
}

func (hs *HttpServer) Stop() error {
	return hs.Server.Shutdown(context.Background())
}

func CreateTLSConfigMinTls(minTlsVersion uint16) *tls.Config {
	conf := &tls.Config{
		// MinVersion: tls.VersionTLS12,
		// MinVersion: tls.VersionTLS11,
		// MinVersion:               tls.VersionTLS10, // weak, only for xp
		MinVersion:               minTlsVersion, // weak, only for xp
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// tls.TLS_RSA_WITH_RC4_128_SHA,
			// tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			// tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,

			// TLS 1.0 - 1.2 cipher suites.
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,

			// TLS 1.3 cipher suites.
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
	return conf
}
