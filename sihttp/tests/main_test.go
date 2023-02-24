package sihttp_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-wonk/si/sihttp"
	_ "github.com/lib/pq"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	longtest, _   = strconv.ParseBool(os.Getenv("LONG_TEST"))

	standardClient *http.Client
	httpServer     *sihttp.Server
	serverAddr     = ":59111"
	remoteAddr     = "http://127.0.0.1:59111"
)

func openClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	dialer := &net.Dialer{Timeout: 5 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       300,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    tlsConfig,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	return sihttp.NewStandardClient(time.Duration(30), tr)
}

func setup() error {
	if onlinetest {
		standardClient = openClient()

		router := gin.Default()
		router.GET("/test/hello", func(c *gin.Context) {
			c.Writer.Write([]byte("hello"))
		})
		router.POST("/test/echo", func(c *gin.Context) {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.Writer.Write([]byte(err.Error()))
				return
			}
			c.Writer.Write(b)
		})

		tlsConfig := sihttp.CreateTLSConfigMinTls(tls.VersionTLS12)
		httpServer = sihttp.NewServer(router, tlsConfig, serverAddr,
			15*time.Second, 15*time.Second)

		go func() {
			if err := httpServer.Start(); err != nil {
				fmt.Println(err)
			}
		}()
	}

	return nil
}

func shutdown() {
	if httpServer != nil {
		httpServer.Stop()
	}
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
}
