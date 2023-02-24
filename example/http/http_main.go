package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/sihttp"
)

func main() {

	dialer := &net.Dialer{Timeout: 3 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	c := &http.Client{
		Timeout:   time.Duration(15) * time.Second,
		Transport: tr,
	}

	client := sihttp.NewClient(c)
	res, err := client.Post("http://127.0.0.1:8080/test/echo", nil, "hello")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(res))

	client2 := sihttp.NewClient(c,
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()))
	type Student struct {
		ID           int    `json:"id"`
		EmailAddress string `json:"email_address"`
		Name         string `json:"name"`
		Borrowed     bool   `json:"borrowed"`
	}

	s := Student{
		ID:           1,
		EmailAddress: "wonk@wonk.org",
		Name:         "wonk",
		Borrowed:     true,
	}
	res, err = client2.Post("http://127.0.0.1:8080/test/echo", nil, &s)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(res))

	var resStudent Student
	err = client2.PostDecode("http://127.0.0.1:8080/test/echo", nil, &s, &resStudent)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resStudent)
}
