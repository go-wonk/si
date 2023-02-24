package sifile_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-wonk/si/v2/sifile"
	"github.com/go-wonk/si/v2/siutils"
)

func TestFile_ReadFrom(t *testing.T) {
	data := "test data to write.\n"
	dataReader := strings.NewReader(data)
	fileName := "./data/TestFile_ReadFrom.txt"

	var fileMode os.FileMode
	fi, err := os.Stat(fileName)
	if err != nil {
		fileMode = 0755
	} else {
		fileMode = fi.Mode()
	}
	f, err := sifile.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, fileMode)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	n, err := f.ReadFrom(dataReader)
	siutils.AssertNilFail(t, err)

	fmt.Println(n)
}
