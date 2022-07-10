package sitcp

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	_tcpConnectionMap sync.Map
)

type TcpConnectionPool interface {
	Get() interface{}
	Put(v interface{})
}

func DeleteTcpConnectionPool(addr string) {
	p, loaded := _tcpConnectionMap.LoadAndDelete(addr)
	if loaded {
		prevSyncPool := p.(*sync.Pool)
		prevSyncPool.New = func() interface{} {
			return nil
		}

		tolerance := 0
		cnt := 0
		for {
			conn := prevSyncPool.Get()
			if conn == nil {
				tolerance++
				if tolerance > 100 {
					break
				}
				continue
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
	err := conn.Reset(opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func PutConn(addr string, c *Conn) {
	getTcpConnectionPool(addr, 1*time.Second).Put(c)
}
