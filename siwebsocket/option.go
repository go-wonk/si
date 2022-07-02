package siwebsocket

type WebsocketOption interface {
	apply(c *Conn)
}

type WebsocketOptionFunc func(c *Conn)

func (o WebsocketOptionFunc) apply(c *Conn) {
	o(c)
}

func WithMessageHandler(h MessageHandler) WebsocketOptionFunc {
	return WebsocketOptionFunc(func(c *Conn) {
		c.SetMessageHandler(h)
	})
}
