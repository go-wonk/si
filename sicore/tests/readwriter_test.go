package sicore_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/go-wonk/si/sicore"
	"github.com/stretchr/testify/assert"
)

func TestReadWriter_Tcp_WriteAndRead(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":10000", 6*time.Second)
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

	s := sicore.NewReadWriterWithValidator(conn, conn, tcpValidator())
	received, err := s.WriteAndRead(createDataToSend())
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	l, err := strconv.Atoi(string(received[:7]))
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Equal(t, l, len(received))
}

func TestReadWriter_Tcp_WriteAndRead2(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":10000", 6*time.Second)
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

	s := sicore.NewReadWriterWithValidator(conn, conn, tcpValidator())
	for i := 0; i < 2; i++ {
		received, err := s.WriteAndRead(createSmallDataToSend())
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		l, err := strconv.Atoi(string(received[:7]))
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, l, len(received))
	}
}

func Test_Basic_Tcp(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	conn, err := net.DialTimeout("tcp", ":10000", 6*time.Second)
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
