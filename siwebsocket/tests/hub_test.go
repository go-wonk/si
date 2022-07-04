package siwebsocket_test

import (
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwebsocket"
)

func TestHub(t *testing.T) {
	hub := siwebsocket.NewHub()
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}

	for i := 0; i < 500; i++ {
		log.Println("connect")
		conn, _, err := siwebsocket.DefaultConn(u, nil)
		siutils.AssertNilFail(t, err)
		c := siwebsocket.NewClientConfiguredWithHub(conn, 10*time.Second, 60*time.Second, 1024000, true, hub,
			siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}))
		go c.Start()

		// c.SetID("9099909")
		err = hub.AddClient(c)
		if err != nil {
			log.Println(err)
			return
		}
	}

	err := hub.Broadcast([]byte("hey"))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(6 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	log.Println("stopped")
	hub.Wait()
	time.Sleep(6 * time.Second)
}
