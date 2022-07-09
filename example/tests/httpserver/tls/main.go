//go:build ignore
// +build ignore

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	ip   string
	port int
	key  string
	cert string
	addr string

	dump bool
	// filePath string
)

func init() {

	flag.StringVar(&ip, "i", "", "ip")
	flag.IntVar(&port, "p", 8081, "port")
	flag.StringVar(&cert, "cert", "./certs/cert.pem", "certificate file path")
	flag.StringVar(&key, "key", "./certs/key.pem", "key file path")
	flag.BoolVar(&dump, "dump", false, "dump request")

	flag.Parse()

	addr = fmt.Sprintf("%v:%v", ip, strconv.Itoa(port))
}
func main() {

	// 라우터, gorilla mux를 쓴다
	router := mux.NewRouter()

	router.HandleFunc("/test/hello", HandleBasic)
	router.HandleFunc("/test/echo", HandleEcho).Methods(http.MethodPost)

	tlsConf := &tls.Config{
		// MinVersion: tls.VersionTLS12,
		// MinVersion: tls.VersionTLS11,
		MinVersion:               tls.VersionTLS10, // weak, only for xp
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
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

	// http 서버 생성
	httpServer := &http.Server{
		Addr:         addr,             // listen 할 주소(ip:port)
		WriteTimeout: 30 * time.Second, // 서버 > 클라이언트 응답
		ReadTimeout:  30 * time.Second, // 클라이언트 > 서버 요청
		Handler:      router,           // mux다
		TLSConfig:    tlsConf,
	}

	log.Fatal(httpServer.ListenAndServeTLS(cert, key))
}

func HandleBasic(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(b))
	w.Write([]byte("hello"))
}

func HandleEcho(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(body))
	w.Write(body)
}
