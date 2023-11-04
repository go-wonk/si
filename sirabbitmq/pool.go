package sirabbitmq

// type UnsafeChannelPool struct {
// 	poolSize int
// 	pool     chan *UnsafeChannel
// }

// func NewUnsafeChannelPool(size int, conn *Conn) *UnsafeChannelPool {
// 	pool := make(chan *UnsafeChannel, size)
// 	for i := 0; i < size; i++ {
// 		pool <- NewUnsafeChannel(conn)
// 	}

// 	p := &UnsafeChannelPool{
// 		poolSize: size,
// 		pool:     pool,
// 	}

// 	return p
// }

// func (p *UnsafeChannelPool) Get() *UnsafeChannel {
// 	return <-p.pool
// }
// func (p *UnsafeChannelPool) Put(c *UnsafeChannel) {
// 	p.pool <- c
// }

// type ChannelPool struct {
// 	poolSize int
// 	pool     chan *Channel
// }

// func NewChannelPool(size int, conn *Conn) *ChannelPool {
// 	pool := make(chan *Channel, size)
// 	for i := 0; i < size; i++ {
// 		pool <- NewChannel(conn)
// 	}

// 	p := &ChannelPool{
// 		poolSize: size,
// 		pool:     pool,
// 	}

// 	return p
// }

// func (p *ChannelPool) Get() *Channel {
// 	return <-p.pool
// }
// func (p *ChannelPool) Put(c *Channel) {
// 	p.pool <- c
// }

type ConnPool struct {
	poolSize int
	pool     chan *Conn
}

func NewConnPool(size int, addr string, prefetch int) *ConnPool {
	pool := make(chan *Conn, size)
	for i := 0; i < size; i++ {
		pool <- NewConn(addr, prefetch)
	}

	p := &ConnPool{
		poolSize: size,
		pool:     pool,
	}

	return p
}

func (p *ConnPool) Get() *Conn {
	return <-p.pool
}
func (p *ConnPool) Put(c *Conn) {
	p.pool <- c
}
func (p *ConnPool) Size() int {
	return p.poolSize
}
func (p *ConnPool) Close() []error {
	errs := make([]error, 0)
	for i := 0; i < p.poolSize; i++ {
		c := p.Get()
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

type UnsafeChannelPool struct {
	poolSize int
	pool     chan *UnsafeChannel
	connPool *ConnPool
}

func NewUnsafeChannelPool(size int, connPool *ConnPool) *UnsafeChannelPool {
	pool := make(chan *UnsafeChannel, size*connPool.Size())
	for i := 0; i < connPool.Size(); i++ {
		conn := connPool.Get()
		for j := 0; j < size; j++ {
			pool <- NewUnsafeChannel(conn)
		}
		connPool.Put(conn)
	}

	p := &UnsafeChannelPool{
		poolSize: size,
		pool:     pool,
		connPool: connPool,
	}

	return p
}

func (p *UnsafeChannelPool) Get() *UnsafeChannel {
	return <-p.pool
}
func (p *UnsafeChannelPool) Put(c *UnsafeChannel) {
	p.pool <- c
}
func (p *UnsafeChannelPool) Close() error {
	p.connPool.Close()
	return nil
}

type ChannelPool struct {
	poolSize int
	pool     chan *Channel
	connPool *ConnPool
}

func NewChannelPool(size int, connPool *ConnPool) *ChannelPool {
	pool := make(chan *Channel, size*connPool.Size())
	for i := 0; i < connPool.Size(); i++ {
		conn := connPool.Get()
		for j := 0; j < size; j++ {
			pool <- NewChannel(conn)
		}
		connPool.Put(conn)
	}

	p := &ChannelPool{
		poolSize: size,
		pool:     pool,
		connPool: connPool,
	}

	return p
}

func (p *ChannelPool) Get() *Channel {
	return <-p.pool
}
func (p *ChannelPool) Put(c *Channel) {
	p.pool <- c
}
func (p *ChannelPool) Close() error {
	p.connPool.Close()
	return nil
}

// import (
// 	"sync"
// )

// var (
// 	_connMap sync.Map
// )

// type ConnPool interface {
// 	Get() interface{}
// 	Put(v interface{})
// }

// func getConnPool(addr string) ConnPool {
// 	p, _ := _connMap.LoadOrStore(addr, &sync.Pool{})
// 	return p.(ConnPool)
// }

// func GetConn(addr string) *Conn {
// 	p := getConnPool(addr)
// 	g := p.Get()
// 	var c *Conn
// 	if g == nil {
// 		Debug("creating a new connection")
// 		c = NewConn(addr)
// 	} else {
// 		c = g.(*Conn)
// 	}
// 	return c
// }

// func PutConn(c *Conn) {
// 	getConnPool(c.addr).Put(c)
// }

// func CloseConns() {
// 	_connMap.Range(func(key, value interface{}) bool {
// 		if mm, ok := value.(*sync.Pool); ok {
// 			for {
// 				g := mm.Get()
// 				if g == nil {
// 					break
// 				}
// 				g.(*Conn).Close()
// 			}
// 			return true
// 		}
// 		return true
// 	})
// }

// var (
// 	_channelMap       sync.Map
// 	_unsafeChannelMap sync.Map
// )

// type ChannelPool interface {
// 	Get() interface{}
// 	Put(v interface{})
// }

// func getChannelPool(conn *Conn) ChannelPool {
// 	p, _ := _channelMap.LoadOrStore(conn.id, &sync.Pool{})

// 	return p.(ChannelPool)
// }

// func GetChannel(conn *Conn) *Channel {
// 	p := getChannelPool(conn)

// 	g := p.Get()

// 	var c *Channel
// 	if g == nil {
// 		Debug("creating a new channel")
// 		c = NewChannel(conn)
// 	} else {
// 		c = g.(*Channel)
// 	}

// 	return c
// }

// func PutChannel(c *Channel) {
// 	getChannelPool(c.conn).Put(c)
// }

// func CloseChannels() {
// 	_channelMap.Range(func(key, value interface{}) bool {
// 		if mm, ok := value.(*sync.Pool); ok {
// 			for {
// 				g := mm.Get()
// 				if g == nil {
// 					break
// 				}
// 				g.(*Channel).Close()
// 			}
// 			return true
// 		}
// 		return true
// 	})
// }

// type UnsafeChannelPool interface {
// 	Get() interface{}
// 	Put(v interface{})
// }

// func getUnsafeChannelPool(conn *Conn) UnsafeChannelPool {
// 	p, _ := _unsafeChannelMap.LoadOrStore(conn.id, &sync.Pool{})

// 	return p.(UnsafeChannelPool)
// }

// func GetUnsafeChannel(conn *Conn) *UnsafeChannel {
// 	p := getUnsafeChannelPool(conn)

// 	g := p.Get()
// 	var c *UnsafeChannel
// 	if g == nil {
// 		Debug("creating a new unsafe channel")
// 		c = NewUnsafeChannel(conn)
// 	} else {
// 		c = g.(*UnsafeChannel)
// 	}

// 	return c
// }

// func PutUnsafeChannel(c *UnsafeChannel) {
// 	getUnsafeChannelPool(c.conn).Put(c)
// }

// func CloseUnsafeChannels() {
// 	_unsafeChannelMap.Range(func(key, value interface{}) bool {
// 		if mm, ok := value.(*sync.Pool); ok {
// 			for {
// 				g := mm.Get()
// 				if g == nil {
// 					break
// 				}
// 				g.(*UnsafeChannel).Close()
// 			}
// 			return true
// 		}
// 		return true
// 	})
// }
