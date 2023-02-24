package sicore_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-wonk/si/v2/sicore"
	"github.com/go-wonk/si/v2/siutils"
	"github.com/stretchr/testify/assert"
)

func Test_Basic_Tcp(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":9999", 6*time.Second)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer conn.Close()

	// tcpConn := conn.(*net.TCPConn)
	// addr, _ := net.ResolveTCPAddr("tcp4", ":10000")
	// conn, err := net.DialTCP("tcp", nil, addr)
	// if !assert.Nil(t, err) {
	// 	t.FailNow()
	// }
	// defer conn.Close()

	err = conn.SetWriteDeadline(time.Now().Add(6 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.SetReadDeadline(time.Now().Add(12 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	err = conn.(*net.TCPConn).SetWriteBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.(*net.TCPConn).SetReadBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	buf := make([]byte, 1024)
	conn.Write(createSmallDataToSend())
	conn.Read(buf)
}

func TestReader_Writer_Tcp_WriteRead(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":9999", 6*time.Second)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer conn.Close()

	// tcpConn := conn.(*net.TCPConn)
	// addr, _ := net.ResolveTCPAddr("tcp4", ":10000")
	// conn, err := net.DialTCP("tcp", nil, addr)
	// if !assert.Nil(t, err) {
	// 	t.FailNow()
	// }
	// defer conn.Close()

	err = conn.SetWriteDeadline(time.Now().Add(6 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.SetReadDeadline(time.Now().Add(12 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	err = conn.(*net.TCPConn).SetWriteBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.(*net.TCPConn).SetReadBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	r := sicore.GetReader(conn, SetTcpEofChecker())
	w := sicore.GetWriter(conn)

	_, err = w.WriteFlush(createDataToSend())
	siutils.AssertNilFail(t, err)

	received, err := r.ReadAll()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	l, err := strconv.Atoi(string(received[:7]))
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, l, len(received))
}

func TestReadWriter_Tcp_Request(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":9999", 6*time.Second)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer conn.Close()

	err = conn.SetWriteDeadline(time.Now().Add(6 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.SetReadDeadline(time.Now().Add(12 * time.Second))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	err = conn.(*net.TCPConn).SetWriteBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	err = conn.(*net.TCPConn).SetReadBuffer(4096)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	// r := sicore.GetReader(conn, SetTcpEofChecker())
	// w := sicore.GetWriter(conn)
	// rw := sicore.NewReadWriter(r, w)

	// rw := sicore.GetReadWriterWithOptions(conn, []sicore.ReaderOption{SetTcpEofChecker()}, conn, nil)
	rw := sicore.GetReadWriterWithReadWriter(conn)
	rw.Reader.ApplyOptions(SetTcpEofChecker())

	recv, err := rw.Request(createDataToSend())
	siutils.AssertNilFail(t, err)

	l, err := strconv.Atoi(string(recv[:7]))
	siutils.AssertNilFail(t, err)
	assert.Equal(t, l, len(recv))
	sicore.PutReadWriter(rw)

	// rw = sicore.GetReadWriterWithOptions(conn, []sicore.ReaderOption{SetTcpEofChecker()}, conn, nil)
	rw = sicore.GetReadWriterWithReadWriter(conn)
	rw.Reader.ApplyOptions(SetTcpEofChecker())
	recv, err = rw.Request(createSmallDataToSend())
	siutils.AssertNilFail(t, err)

	l, err = strconv.Atoi(string(recv[:7]))
	siutils.AssertNilFail(t, err)
	assert.Equal(t, l, len(recv))
	sicore.PutReadWriter(rw)
}

func createDataToSend() []byte {

	dataToSend := strings.Repeat("a", 900000)
	dataLength := len(dataToSend) + 7
	dataLengthStr := fmt.Sprintf("%07d", dataLength)
	return []byte(dataLengthStr + dataToSend)
}

func createSmallDataToSend() []byte {
	dataToSend := strings.Repeat("a", 10)
	dataLength := len(dataToSend) + 7
	dataLengthStr := fmt.Sprintf("%07d", dataLength)
	return []byte(dataLengthStr + dataToSend)
}

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

func SetTcpEofChecker() sicore.ReaderOption {
	return sicore.ReaderOptionFunc(func(r *sicore.Reader) {
		r.SetEofChecker(&TcpEOFChecker{})
	})
}

// func tcpValidator() sicore.ReadValidator {
// 	return sicore.ValidateFunc(func(b []byte, errIn error) (bool, error) {
// 		if errIn == nil || errIn == io.EOF {
// 			lenStr := string(b[:7])
// 			lenProt, err := strconv.ParseInt(lenStr, 10, 64)
// 			if err != nil {
// 				return false, errors.New("cannot find data length")
// 			}

// 			receivedAll := int(lenProt) == len(b)
// 			if receivedAll {
// 				return true, nil
// 			}

// 			if errIn == io.EOF {
// 				return false, errors.New("not received all but EOF")
// 			}
// 			return false, nil
// 		}

// 		return false, errIn
// 	})
// }
