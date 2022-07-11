package sitcp

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	lock              sync.RWMutex
	errLimit          int64 = 10
	errCount          int64
	_tcpConnectionMap sync.Map
)

type TcpConnectionPool interface {
	Get() interface{}
	Put(v interface{})
}

func deleteTcpConnectionPool(addr string) {
	p, loaded := _tcpConnectionMap.LoadAndDelete(addr)
	if loaded {
		prevSyncPool := p.(*sync.Pool)
		prevSyncPool.New = func() interface{} {
			return nil
		}

		cnt := 0
		for {
			conn := prevSyncPool.Get()
			if conn == nil {
				break
			}
			cnt++
			log.Println("delete: closed", conn.(*Conn).LocalAddr())
			conn.(*Conn).Close()
		}
		log.Println(cnt)
	}
}

func getTcpConnectionPool(addr string, timeout time.Duration, opts ...TcpOption) TcpConnectionPool {
	p, _ := _tcpConnectionMap.LoadOrStore(addr, &sync.Pool{
		New: func() interface{} {
			log.Println("new connection")
			c, err := DialTimeout(addr, timeout, opts...)
			if err != nil {
				return nil
			}
			return c
		},
	})

	return p.(TcpConnectionPool)
}

func GetConn(addr string, timeout time.Duration, opts ...TcpOption) (*Conn, error) {
	c := getTcpConnectionPool(addr, timeout, opts...).Get()
	if c == nil {
		return nil, errors.New("cannot create tcp connection")
	}
	conn := c.(*Conn)
	err := conn.reset(opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// PutConn puts connection back to the pool with key, addr.
func PutConn(addr string, c *Conn) {
	if c.Err() != nil {
		c.Close()

		lock.Lock()
		errCount++
		if errCount > errLimit {
			deleteTcpConnectionPool(addr)
			errCount = 0
		}
		lock.Unlock()
		return
	}
	getTcpConnectionPool(addr, time.Second).Put(c)
}
