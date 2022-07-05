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

	done       chan struct{}
	stop       chan string
	clientDone chan struct{}
	terminated chan struct{}

	clientStorage ClientStorage
}

func NewHub() *Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		clients:       sync.Map{},
		done:          make(chan struct{}),
		stop:          make(chan string, 1),
		clientDone:    make(chan struct{}),
		terminated:    make(chan struct{}),
		clientStorage: &NopRouteStorage{},
	}
	go h.waitStop()

	return h
}

func (h *Hub) Run() {
	// defer func() {
	// 	close(h.terminated)
	// }()
	for {
		select {
		case <-h.done:
			return
		case client := <-h.register:
			// h.clients[client] = true
			loadedClient, exist := h.clients.LoadOrStore(client.id, client)
			if exist {
				// Stop will lead to removeClient method called.
				loadedClient.(*Client).Stop()
				h.clients.Store(client.id, client)
			}
		case client := <-h.unregister:
			deletedClient, exist := h.clients.LoadAndDelete(client.id)
			if exist {
				log.Println("deleted:", deletedClient.(*Client).id)
			} else {
				log.Println("not found client:", client.id)
			}
			// client.Stop()
			// client.Wait()
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
	<-h.stop             // wait until Stop method is called
	close(h.clientDone)  // prevent from sending into register/unregister/broadcast channel
	close(h.done)        // stops Run method
	h.removeAllClients() // stops and closes all clients and remove from clients map
	close(h.terminated)
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

func (h *Hub) addClient(client *Client) error {
	select {
	case <-h.clientDone:
		return errors.New("register buffer closed")
	case h.register <- client:
		return h.clientStorage.Store(client.id, "", "", "", "")
	}
}

func (h *Hub) removeClient(client *Client) error {
	select {
	case <-h.clientDone:
		return errors.New("register buffer closed")
	case h.unregister <- client:
		return h.clientStorage.Delete(client.id)
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
		err := value.(*Client).Stop()
		if err != nil {
			log.Println("removeAllClients err:", err)
		}
		value.(*Client).Wait()
		h.clients.Delete(value.(*Client).id)
		h.clientStorage.Delete(value.(*Client).id)
		return true
	})

	return nil
}

func (h *Hub) RemoveRandomClient() error {
	h.clients.Range(func(key interface{}, value interface{}) bool {
		err := value.(*Client).Stop()
		if err != nil {
			log.Println("RemoveRandomClient err:", err)
		}
		// value.(*Client).Wait()
		return false
	})

	return nil
}

func (h *Hub) LenClients() int {
	lenClients := 0
	h.clients.Range(func(key interface{}, value interface{}) bool {
		lenClients++
		log.Println("clients left in hub", value.(*Client).id)
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
