package sitcp_test

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sitcp"
	"github.com/go-wonk/si/siutils"
)

type TcpEOFChecker struct{}

func (c TcpEOFChecker) Check(b []byte, errIn error) (bool, error) {
	if errIn == nil || errIn == io.EOF {
		lenStr := string(b[:7])
		lenProt, err := strconv.ParseInt(lenStr, 10, 64)
		if err != nil {
			return false, errors.New("cannot find data length")
		}

		receivedAll := int(lenProt) == len(b)
		if receivedAll {
			return true, nil
		}

		if errIn == io.EOF {
			return false, errors.New("not received all but EOF")
		}
		return false, nil
	}

	return false, errIn
}
func createSmallDataToSend() []byte {
	dataToSend := strings.Repeat("a", 10)
	dataLength := len(dataToSend) + 7
	dataLengthStr := fmt.Sprintf("%07d", dataLength)
	return []byte(dataLengthStr + dataToSend)
}
func TestConn_Request(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	defaultConn, err := sitcp.DefaultTcpConn("127.0.0.1:10000", 3, 3, 3, 4096, 4096)
	siutils.AssertNilFail(t, err)
	defer defaultConn.Close()

	conn := sitcp.NewConn(defaultConn)
	conn.SetReaderOption(sicore.SetEofChecker(&TcpEOFChecker{}))

	res, err := conn.Request(createSmallDataToSend())
	siutils.AssertNilFail(t, err)

	fmt.Println(res)
}
