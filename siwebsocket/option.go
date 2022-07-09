package siwebsocket

import "github.com/go-wonk/si/sicore"

// ClientOption is an interface with apply method.
type ClientOption interface {
	apply(c *Client)
}

// ClientOptionFunc wraps a function to conforms to ClientOption interface
type ClientOptionFunc func(c *Client)

// apply implements ClientOption's apply method.
func (o ClientOptionFunc) apply(c *Client) {
	o(c)
}

// WithMessageHandler sets h to c.
func WithMessageHandler(h MessageHandler) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetMessageHandler(h)
	})
}

// WithReaderOpt sets ro to c.
func WithReaderOpt(ro sicore.ReaderOption) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.appendReaderOpt(ro)
	})
}

// WithID sets id to c's ID.
func WithID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetID(id)
	})
}

// WithUserID sets id to c's userID
func WithUserID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetUserID(id)
	})
}

// WithUserGroupID sets id to c's userGroupID
func WithUserGroupID(id string) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetUserGroupID(id)
	})
}

// HubOption is an interface with apply method.
type HubOption interface {
	apply(h *Hub)
}

// HubOptionFunc wraps a function to conforms to HubOption interface.
type HubOptionFunc func(h *Hub)

// apply implements HubOption's apply method.
func (o HubOptionFunc) apply(h *Hub) {
	o(h)
}

// WithRouter sets r to h's router.
func WithRouter(r Router) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.SetRouter(r)
	})
}

// WithHubAddr sets addr to h's hubAddr.
func WithHubAddr(addr string) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.SetHubAddr(addr)
	})
}

// WithHubPath sets path to h's hubPath.
func WithHubPath(path string) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.SetHubPath(path)
	})
}

// WithAfterDeleteClient sets f to h's afterDeleteClient.
func WithAfterDeleteClient(f func(*Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.afterDeleteClient = f
	})
}

// WithAfterStoreClient sets f to h's afterStoreClient.
func WithAfterStoreClient(f func(*Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.afterStoreClient = f
	})
}
