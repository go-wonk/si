package sitcp_test

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sitcp"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

type TcpEOFChecker struct{}

func (c TcpEOFChecker) Check(b []byte, errIn error) (bool, error) {
	if errIn == nil || errIn == io.EOF {
		if len(b) == 0 {
			return false, errors.New("conn may have been closed")
		}
		lenStr := string(b[:7])
		lenProt, err := strconv.ParseInt(lenStr, 10, 64)
		if err != nil {
			return false, errors.New("cannot find data length " + string(b))
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
func createSmallDataToSend(i, j int) []byte {
	istr := strconv.Itoa(i)
	jstr := strconv.Itoa(j)
	dataToSend := strings.Repeat(istr+jstr, 10)
	dataLength := len(dataToSend) + 7
	dataLengthStr := fmt.Sprintf("%07d", dataLength)
	return []byte(dataLengthStr + dataToSend)
}
func TestConn_Request(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	conn, err := sitcp.DialTimeout("127.0.0.1:9999", 3*time.Second,
		sitcp.WithReaderOpt(sicore.SetEofChecker(&TcpEOFChecker{})))
	siutils.AssertNilFail(t, err)
	defer conn.Close()

	siutils.AssertNilFail(t, err)
	res, err := conn.Request(createSmallDataToSend(1, 2))
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))
}

func TestConn_Request_Concurrent(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			conn, err := sitcp.DialTimeout("127.0.0.1:9999", 3*time.Second,
				sitcp.WithEofChecker(&TcpEOFChecker{}),
				sitcp.WithWriteTimeout(3*time.Second),
				sitcp.WithReadTimeout(3*time.Second),
				sitcp.WithWriteBufferSize(4096),
				sitcp.WithReadBufferSize(4096),
			)
			if !assert.Nil(t, err) {
				return
			}

			for j := 0; j < 100; j++ {

				sendData := createSmallDataToSend(i, j)
				res, err := conn.Request(sendData)
				siutils.AssertNilFail(t, err)
				assert.EqualValues(t, sendData, string(res))
				log.Println(string(res))
			}

			conn.Close()
		}(i)
	}
	wg.Wait()
}

func TestConnPool_Request(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	addr := "127.0.0.1:10000"
	for i := 0; i < 50; i++ {
		c, _ := sitcp.DialTimeout(addr, 3*time.Second,
			sitcp.WithEofChecker(&TcpEOFChecker{}))
		log.Println(c.LocalAddr())
		sitcp.PutConn(addr, c)
	}

	var val time.Duration = 50
	for i := 0; i < 500; i++ {
		conn, err := sitcp.GetConn(addr, 3*time.Second,
			sitcp.WithEofChecker(&TcpEOFChecker{}),
			sitcp.WithWriteTimeout(3*time.Second),
			sitcp.WithReadTimeout(3*time.Second))
		// siutils.AssertNilFail(t, err)
		if err != nil {
			log.Println("conn:", err)
			time.Sleep(val * time.Millisecond)
			continue
		}

		sendData := createSmallDataToSend(1, 2)
		res, err := conn.Request(sendData)
		// siutils.AssertNilFail(t, err)
		if err != nil {
			sitcp.PutConn(addr, conn)
			log.Println("request:", err)
			time.Sleep(val * time.Millisecond)
			continue
		}

		fmt.Println(conn.LocalAddr(), string(res))
		assert.EqualValues(t, sendData, string(res))
		sitcp.PutConn(addr, conn)

		time.Sleep(val * time.Millisecond)

		// if i > 50 {
		// 	val = 3000
		// }
	}
}

func TestConnPool_Request_Concurrent(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	addr := "127.0.0.1:10000"
	for i := 0; i < 50; i++ {
		c, _ := sitcp.DialTimeout(addr, 3*time.Second,
			sitcp.WithEofChecker(&TcpEOFChecker{}))
		log.Println(c.LocalAddr())
		sitcp.PutConn(addr, c)
	}

	wg := &sync.WaitGroup{}
	var val time.Duration = 100
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				conn, err := sitcp.GetConn(addr, 3*time.Second,
					sitcp.WithEofChecker(&TcpEOFChecker{}),
					sitcp.WithWriteTimeout(3*time.Second),
					sitcp.WithReadTimeout(3*time.Second))
				// siutils.AssertNilFail(t, err)
				if err != nil {
					log.Println("conn:", err)
					continue
				}

				sendData := createSmallDataToSend(i, j)
				res, err := conn.Request(sendData)
				// siutils.AssertNilFail(t, err)
				if err != nil {
					sitcp.PutConn(addr, conn)
					log.Println("request:", err)
					continue
				}

				fmt.Println(conn.LocalAddr(), string(res))
				assert.EqualValues(t, sendData, string(res))
				sitcp.PutConn(addr, conn)

				time.Sleep(val * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}
