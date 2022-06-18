package sifile

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-wonk/si/siutils"
)

func TestNewFile(t *testing.T) {
	f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0777)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	byt, err := f.ReadAll()
	siutils.AssertNilFail(t, err)
	fmt.Println(string(byt) + "1")

	_, err = f.WriteFlush([]byte("hey\n"))
	siutils.AssertNilFail(t, err)

	byt, err = f.ReadAllFrom(0)
	siutils.AssertNilFail(t, err)
	fmt.Println(string(byt) + "2")

	_, err = f.WriteFlush([]byte("hey2\n"))
	siutils.AssertNilFail(t, err)

	byt, err = f.ReadAllFrom(0)
	siutils.AssertNilFail(t, err)
	fmt.Println(string(byt) + "3")

	_, err = f.WriteFlush([]byte("hey3\n"))
	siutils.AssertNilFail(t, err)

	byt, err = f.ReadAllFrom(0)
	siutils.AssertNilFail(t, err)
	fmt.Println(string(byt) + "4")
}
