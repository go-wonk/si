package siwebsocket_test

import (
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwebsocket"
)

func TestWebsocket(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/echo"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	defer conn.Close()

	siconn := siwebsocket.NewConn(conn)
	go siconn.ReadPump()
	for i := 0; i < 20; i++ {
		go func(i int) {
			defer func() {
				log.Println("returning", i)
			}()
			for {
				// time.Sleep(100 * time.Millisecond)

				rn := rand.Intn(1000)
				if err := siconn.Send([]byte(strconv.Itoa(rn))); err != nil {
					if err != siwebsocket.ErrDataChannelClosed {
						log.Println(err)
					}
					return
				}
			}

		}(i)
	}

	// go func() {
	time.Sleep(6 * time.Second)
	siconn.Stop()
	siconn.Wait()
	log.Println("terminated")
	// siconn.CloseGracefully()
	// }()

}

func TestWebsocket2(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/echo/randomclose"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	defer conn.Close()

	siconn := siwebsocket.NewConn(conn)
	go siconn.ReadPump()
	for i := 0; i < 20; i++ {
		go func(i int) {
			defer func() {
				log.Println("returning", i)
			}()
			for {
				// time.Sleep(100 * time.Millisecond)

				rn := rand.Intn(1000)
				// rn = 1001
				if err := siconn.Send([]byte(strconv.Itoa(rn))); err != nil {
					if err != siwebsocket.ErrDataChannelClosed {
						log.Println(err)
					}
					return
				}
			}

		}(i)
	}

	siconn.Wait()
	log.Println("terminated")

}

func TestWebsocket3(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	c := siwebsocket.NewConn(conn)
	go c.ReadPump()
	time.Sleep(12 * time.Second)
	c.Stop()
	c.Wait()
	log.Println("terminated 1")

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		u2 := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}
		conn2, _, err := siwebsocket.DefaultConn(u2, nil)
		siutils.AssertNilFail(t, err)
		c2 := siwebsocket.NewConn(conn2)
		go c2.ReadPump()
		time.Sleep(12 * time.Second)
		c2.Stop()
		c2.Wait()
		log.Println("terminated 2")
	}

}

func TestWebsocket4(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "192.168.0.92:48080", Path: "/idle"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	c := siwebsocket.NewConn(conn)
	go c.ReadPump()
	time.Sleep(12000 * time.Second)
	c.Stop()
	c.Wait()
	log.Println("terminated 1")

}

func TestWebsocket_EchoIdle(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "192.168.0.92:48080", Path: "/echo"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	c := siwebsocket.NewConn(conn)
	go c.ReadPump()

	c.Wait()
	log.Println("terminated 1")

}

func TestWebsocket_EchoStop(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "192.168.0.92:48080", Path: "/echo"}
	conn, _, err := siwebsocket.DefaultConn(u, nil)
	siutils.AssertNilFail(t, err)
	c := siwebsocket.NewConn(conn)
	go c.ReadPump()

	time.Sleep(10 * time.Second)
	c.Stop()

	c.Wait()
	log.Println("terminated 1")

	// disconnect network right after readPump. Stop results in calling closeMessage.
	// if network is kept disconnected, then "read tcp 192.168.0.12:63300->192.168.0.92:48080: i/o timeout" error occurs.
	// else if network is reconnected, then "websocket: close 1000 (normal)" error occurs.(same result shown when there is no ping from client)
}

func TestWebsocket_Push(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}

	for i := 0; i < 5; i++ {
		log.Println("connect")
		conn, _, err := siwebsocket.DefaultConn(u, nil)
		siutils.AssertNilFail(t, err)
		c := siwebsocket.NewConn(conn, siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}))
		go c.ReadPump()

		time.Sleep(1 * time.Second)
		c.Stop()
		c.Wait()
	}

	for i := 0; i < 3; i++ {
		log.Println("connect")
		conn, _, err := siwebsocket.DefaultConn(u, nil)
		siutils.AssertNilFail(t, err)
		c := siwebsocket.NewConn(conn, siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}))
		go c.ReadPump()

		c.Wait()
		if err := c.ReadErr(); err != nil {
			log.Println(err)
		}
		if err := c.WriteErr(); err != nil {
			log.Println(err)
		}
	}
	log.Println("terminated")
}
