package siwebsocket

import (
	"net/http"
	"time"

	"github.com/go-wonk/si/v2/sicore"
)

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

func WithHub(h sicore.Hub) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.SetHub(h)
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

func WithWriteWait(writeWait time.Duration) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.writeWait = writeWait
	})
}
func WithReadWait(readWait time.Duration) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.readWait = readWait
		c.pingPeriod = (readWait * 9) / 10
	})
}
func WithMaxMessageSize(maxMessageSize int) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.maxMessageSize = maxMessageSize
	})
}
func WithUsePingPong(usePingPong bool) ClientOptionFunc {
	return ClientOptionFunc(func(c *Client) {
		c.usePingPong = usePingPong
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
func WithAfterDeleteClient(f func(sicore.Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.afterDeleteClient = f
	})
}

// WithAfterStoreClient sets f to h's afterStoreClient.
func WithAfterStoreClient(f func(sicore.Client, error)) HubOptionFunc {
	return HubOptionFunc(func(h *Hub) {
		h.afterStoreClient = f
	})
}

// UpgraderOption is an interface with apply method.
type UpgraderOption interface {
	apply(u *upgraderConfig)
}

// UpgraderOptionFunc wraps a function to conforms to ClientOption interface
type UpgraderOptionFunc func(u *upgraderConfig)

// apply implements UpgraderOption's apply method.
func (o UpgraderOptionFunc) apply(u *upgraderConfig) {
	o(u)
}

func WithUpgradeHandshakeTimeout(timeout time.Duration) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.handshakeTimeout = timeout
	})
}

func WithUpgradeReadBufferSize(bufferSize int) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.readBufferSize = bufferSize
	})
}
func WithUpgradeWriteBufferSize(bufferSize int) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.writeBufferSize = bufferSize
	})
}

func WithUpgradeSubprotocols(protocols []string) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.subprotocols = protocols
	})
}

func WithUpgradeError(f func(w http.ResponseWriter, r *http.Request, status int, reason error)) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.errorFunc = f
	})
}

// WithUpgradeCheckOrigin sets f to u's CheckOrigin.
func WithUpgradeCheckOrigin(f func(r *http.Request) bool) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.checkOrigin = f
	})
}

func WithUpgradeEnableCompression(enableCompression bool) UpgraderOptionFunc {
	return UpgraderOptionFunc(func(u *upgraderConfig) {
		u.enableCompression = enableCompression
	})
}
