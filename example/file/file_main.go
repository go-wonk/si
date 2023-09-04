package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/sifile"
)

func main() {
	f, err := sifile.OpenFile("data/test.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	n, err := f.WriteFlush([]byte("hello"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(n)

	f2, err := sifile.OpenFile("data/encode.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm,
		sifile.WithWriterOpt(sicore.SetJsonEncoder()),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f2.Close()
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
	if err := f2.EncodeFlush(&s); err != nil {
		log.Fatal(err)
		return
	}

	f3, err := sifile.OpenFile("data/encode.txt", os.O_CREATE|os.O_RDONLY, os.ModePerm,
		sifile.WithReaderOpt(sicore.SetJsonDecoder()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f3.Close()
	var resStudent Student
	if err := f3.Decode(&resStudent); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resStudent)

	fmt.Println(f3.ReadLine())
}
