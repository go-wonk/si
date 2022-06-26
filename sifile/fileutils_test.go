package sifile_test

import (
	"fmt"
	"testing"

	"github.com/go-wonk/si/sifile"
	"github.com/go-wonk/si/siutils"
)

func TestListAll(t *testing.T) {
	list, err := sifile.ListDir("./tests")
	siutils.AssertNilFail(t, err)

	for _, f := range list {
		fi, err := f.Info()
		siutils.AssertNilFail(t, err)
		fmt.Println(f.Path, f.IsDir(), fi.Size())
	}
}
