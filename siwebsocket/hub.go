package siwebsocket

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients sync.Map
	router  Router

	// channel to broadcast message to connected clients
	broadcast chan []byte

	// channel to add clients to `clients` map
	register chan *Client

	// channel to remove clients from `clients` map
	unregister chan *Client

	runDone    chan struct{}
	stopClient chan string
	clientDone chan struct{}
	terminated chan struct{}

	// Time allowed to write a message to the peer.
	writeWait time.Duration
	// Time allowed to read the next pong message from the peer.
	readWait time.Duration
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration
	// Maximum message size allowed from peer.
	maxMessageSize int
	// use ping/pong
	usePingPong bool

	hubAddr string
	hubPath string

	// handlers
	afterDeleteClient func(c *Client, ok bool)
	afterStoreClient  func(c *Client, ok bool)
}

func NewHub(hubAddr, hubPath string, writeWait time.Duration, readWait time.Duration,
	maxMessageSize int, usePingPong bool, opts ...HubOption) *Hub {

	pingPeriod := (readWait * 9) / 10

	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		clients:    sync.Map{},
		runDone:    make(chan struct{}),
		stopClient: make(chan string, 1),
		clientDone: make(chan struct{}),
		terminated: make(chan struct{}),
		router:     &NopRouter{},

		hubAddr:        hubAddr,
		hubPath:        hubPath,
		writeWait:      writeWait,
		readWait:       readWait,
		pingPeriod:     pingPeriod,
		maxMessageSize: maxMessageSize,
		usePingPong:    usePingPong,
	}

	for _, o := range opts {
		o.apply(h)
	}

	if h.afterDeleteClient == nil {
		h.afterDeleteClient = func(c *Client, ok bool) {
			// nothing
		}
	}

	if h.afterStoreClient == nil {
		h.afterStoreClient = func(c *Client, ok bool) {
			// nothing
		}
	}

	go h.waitStop()

	return h
}

func (h *Hub) CreateAndAddClient(conn *websocket.Conn, opts ...ClientOption) (*Client, error) {

	c := NewClientConfigured(conn, h.writeWait, h.readWait, h.maxMessageSize, h.usePingPong, opts...)
	c.hub = h

	err := h.addClient(c)
	if err != nil {
		c.Stop()
		c.Wait()
		return nil, err
	}
	return c, nil
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.runDone:
			return
		case client := <-h.register:
			loadedClient, exist := h.clients.LoadOrStore(client.id, client)
			if exist {
				// Stop will lead to removeClient method called.
				// Do not call removeClient method here.
				loadedClient.(*Client).Stop()
				h.clients.Store(client.id, client)
			}
			h.afterStoreClient(client, exist)
		case client := <-h.unregister:
			// stopped clients with connection closed are received here.
			// remove them from `clients` map
			_, exist := h.clients.LoadAndDelete(client.id)
			h.afterDeleteClient(client, exist)
		case message := <-h.broadcast:
			// iterating over the map here may cause other channels blocked.
			h.clients.Range(func(key interface{}, value interface{}) bool {
				// time.Sleep(1 * time.Second)
				value.(*Client).Send(message)
				return true
			})
		}
	}
}

func (h *Hub) waitStop() {
	<-h.stopClient       // wait until Stop method is called
	close(h.clientDone)  // prevent from sending into register/unregister/broadcast channel
	close(h.runDone)     // stops Run method
	h.removeAllClients() // stops and closes all clients and remove from clients map
	close(h.terminated)
}

func (h *Hub) Stop() error {
	select {
	case h.stopClient <- "stop":
	default:
		return ErrStopChannelFull
	}
	return nil
}

func (h *Hub) Wait() {
	<-h.terminated
}

func (h *Hub) addClient(client *Client) error {
	select {
	case <-h.clientDone:
		h.afterStoreClient(client, false)
		return errors.New("register buffer closed")
	case h.register <- client:
		return h.router.Store(context.Background(), client.id, client.userID, client.userGroupID, h.hubAddr, h.hubPath)
	}
}

func (h *Hub) removeClient(client *Client) error {
	select {
	case <-h.clientDone:
		h.afterDeleteClient(client, false)
		return errors.New("register buffer closed")
	case h.unregister <- client:
		return h.router.Delete(context.Background(), client.id)
	}
}

func (h *Hub) Broadcast(message []byte) error {
	select {
	case <-h.clientDone:
		return errors.New("register buffer closed")
	case h.broadcast <- message:
	}

	return nil
}

func (h *Hub) removeAllClients() error {
	h.clients.Range(func(key interface{}, value interface{}) bool {
		value.(*Client).Stop()
		value.(*Client).Wait()
		h.clients.Delete(value.(*Client).id)
		h.router.Delete(context.Background(), value.(*Client).id)
		return true
	})

	return nil
}

func (h *Hub) RemoveRandomClient() error {
	h.clients.Range(func(key interface{}, value interface{}) bool {
		value.(*Client).Stop()
		return false
	})

	return nil
}

func (h *Hub) LenClients() int {
	lenClients := 0
	h.clients.Range(func(key interface{}, value interface{}) bool {
		lenClients++
		return true
	})

	return lenClients
}

func (h *Hub) SendMessage(id string, msg []byte) error {
	if c, ok := h.clients.Load(id); !ok {
		return errors.New("client not found, " + id)
	} else {
		err := c.(*Client).Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) SendMessageWithResult(id string, msg []byte) error {
	if c, ok := h.clients.Load(id); !ok {
		return errors.New("client not found, " + id)
	} else {
		m := NewMsg(msg)
		err := c.(*Client).SendMsg(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) SetRouter(r Router) {
	h.router = r
}
