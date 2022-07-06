package siwebsocket

import "github.com/go-wonk/si/sicore"

type ClientOption interface {
	apply(c *Client)
}

type ClientOptionFunc func(c *Client)

func (o ClientOptionFunc) apply(c *Client) {
	o(c)
}

func WithMessageHandler(h MessageHandler) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetMessageHandler(h)
	})
}

func WithReaderOpt(ro sicore.ReaderOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.appendReaderOpt(ro)
	})
}

func WithID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetID(id)
	})
}

func WithUserID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetUserID(id)
	})
}

func WithUserGroupID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetUserGroupID(id)
	})
}

type HubOption interface {
	apply(h *Hub)
}

type HubOptionFunc func(h *Hub)

func (o HubOptionFunc) apply(h *Hub) {
	o(h)
}

func WithRouter(r Router) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.SetRouter(r)
	})
}
