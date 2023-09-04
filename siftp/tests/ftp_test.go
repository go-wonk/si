package siftp_test

import (
	"fmt"
	"testing"

	"github.com/go-wonk/si/v2/siftp"
	"github.com/go-wonk/si/v2/siutils"
)

func TestRead(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	c := siftp.NewClient("", "", "")
	res, err := c.ReadFile("Version.ini")
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))
}

func TestWriteFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	c := siftp.NewClient("", "", "")
	err := c.WriteFile("test.txt", []byte("test upload"))
	siutils.AssertNilFail(t, err)

}
