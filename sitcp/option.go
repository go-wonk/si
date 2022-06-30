package sitcp

type TcpOption interface {
	apply(c *Conn)
}

type TcpOptionFunc func(*Conn)

func (s TcpOptionFunc) apply(c *Conn) {
	s(c)
}

// func SetEOFChecker() TcpOption {
// 	return TcpOptionFunc(func(c *Conn) {
// 		c.SetEOFChecker(&sicore.TcpEOFChecker{})
// 	})
// }
