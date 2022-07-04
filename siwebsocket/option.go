package siwebsocket

import "github.com/go-wonk/si/sicore"

type WebsocketOption interface {
	apply(c *Client)
}

type WebsocketOptionFunc func(c *Client)

func (o WebsocketOptionFunc) apply(c *Client) {
	o(c)
}

func WithMessageHandler(h MessageHandler) WebsocketOptionFunc {
	return WebsocketOptionFunc(func(c *Client) {
		c.SetMessageHandler(h)
	})
}

func WithReaderOpt(ro sicore.ReaderOption) WebsocketOptionFunc {
	return WebsocketOptionFunc(func(c *Client) {
		c.appendReaderOpt(ro)
	})
}
