package siwebsocket

import (
	"errors"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	// clients map[*Client]bool
	clients sync.Map

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	done         chan struct{}
	stop         chan string
	registerDone chan struct{}
	terminated   chan struct{}
}

func NewHub() *Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		clients:      sync.Map{},
		done:         make(chan struct{}),
		stop:         make(chan string, 1),
		registerDone: make(chan struct{}),
		terminated:   make(chan struct{}),
	}
	go h.waitStop()

	return h
}

func (h *Hub) Run() {
	defer func() {
		close(h.terminated)
	}()
	for {
		select {
		case <-h.done:
			return
		case client := <-h.register:
			// h.clients[client] = true
			loadedClient, exist := h.clients.LoadOrStore(client.id, client)
			if exist {
				// h.RemoveClient(loadedClient.(*Client))
				loadedClient.(*Client).Stop()
				// loadedClient.(*Client).Wait()
				h.clients.Store(client.id, client)
			}
		case client := <-h.unregister:
			deletedClient, exist := h.clients.LoadAndDelete(client.id)
			if exist {
				log.Println("deleted:", deletedClient.(*Client).id)
			} else {
				log.Println("not found client:", client.id)
			}
			client.Stop()
			client.Wait()
		case message := <-h.broadcast:
			h.clients.Range(func(key interface{}, value interface{}) bool {
				err := value.(*Client).Send(message)
				if err != nil {
					log.Println("broadcast:", err)
				}
				return true
			})
		}
	}
}

func (h *Hub) waitStop() {
	<-h.stop // wait until Stop method is called
	close(h.registerDone)
	h.CloseAllClients()
	close(h.done)
}

func (h *Hub) Stop() error {
	select {
	case h.stop <- "stop":
	default:
		return ErrStopChannelFull
	}
	return nil
}

func (h *Hub) Wait() {
	<-h.terminated
}

func (h *Hub) AddClient(client *Client) error {
	select {
	case <-h.registerDone:
		return errors.New("register buffer closed")
	case h.register <- client:
	}
	return nil
}

func (h *Hub) removeClient(client *Client) error {
	h.unregister <- client
	return nil
}

func (h *Hub) CloseAllClients() error {
	h.clients.Range(func(key interface{}, value interface{}) bool {
		err := value.(*Client).Stop()
		if err != nil {
			log.Println("CloseAllClients err:", err)
		}
		// value.(*Client).Wait()
		return true
	})

	return nil
}

func (h *Hub) Broadcast(message []byte) error {
	select {
	case <-h.registerDone:
		return errors.New("register buffer closed")
	case h.broadcast <- message:
	}

	return nil
}
