package siwebsocket_test

import (
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwebsocket"
	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	hub := siwebsocket.NewHub("http://127.0.0.1:8080", "/path/_push", 10*time.Second, 60*time.Second, 1024000, true)
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}

	for i := 0; i < 50; i++ {
		log.Println("connect")
		conn, _, err := siwebsocket.DefaultConn(u, nil)
		siutils.AssertNilFail(t, err)

		_, err = siwebsocket.NewClient(conn,
			siwebsocket.WithWriteWait(10*time.Second),
			siwebsocket.WithReadWait(60*time.Second),
			siwebsocket.WithMaxMessageSize(1024000),
			siwebsocket.WithUsePingPong(true),
			siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}),
			siwebsocket.WithHub(hub),
		)
		if err != nil {
			log.Println(err)
			return
		}
	}

	err := hub.Broadcast([]byte("hey"))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(4 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	log.Println("stopped")
	hub.Wait()
	// time.Sleep(6 * time.Second)
}

func TestHub2(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	hub := siwebsocket.NewHub("http://127.0.0.1:8080", "/path/_push", 10*time.Second, 60*time.Second, 1024000, true)
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push"}

	go func() {
		num := 0
		for {
			time.Sleep(80 * time.Millisecond)
			log.Println("connect")
			conn, _, err := siwebsocket.DefaultConn(u, nil)
			if err != nil {
				log.Println(err)
				return
			}

			c, err := siwebsocket.NewClient(conn,
				siwebsocket.WithWriteWait(10*time.Second),
				siwebsocket.WithReadWait(60*time.Second),
				siwebsocket.WithMaxMessageSize(1024000),
				siwebsocket.WithUsePingPong(true),
				siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}),
				siwebsocket.WithHub(hub),
			)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(c.GetID())

			num++
			if num > 300 {
				// hub.RemoveRandomClient()
				num--
				break
			}
		}
	}()

	go func() {
		for {
			time.Sleep(71 * time.Millisecond)
			err := hub.Broadcast([]byte("hey"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	time.Sleep(4 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	hub.Wait()
	// time.Sleep(12 * time.Second)
	log.Println("stopped", hub.LenClients())
}

func test() int {
	hub := siwebsocket.NewHub("http://127.0.0.1:8080", "/path/_push", 10*time.Second, 60*time.Second, 1024000, true)
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push/randomclose"}

	log.Println("start")
	go func() {
		num := 0
		for {
			time.Sleep(24 * time.Millisecond)
			conn, _, err := siwebsocket.DefaultConn(u, nil)
			if err != nil {
				log.Println(err)
				return
			}

			c, err := siwebsocket.NewClient(conn,
				siwebsocket.WithWriteWait(10*time.Second),
				siwebsocket.WithReadWait(60*time.Second),
				siwebsocket.WithMaxMessageSize(1024000),
				siwebsocket.WithUsePingPong(true),
				siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageHandler{}),
				siwebsocket.WithHub(hub),
			)
			if err != nil {
				log.Println(err)
				return
			}

			rn := rand.Intn(1000)
			if rn == 0 {
				hub.RemoveRandomClient()
			}
			go func() {
				for {
					time.Sleep(54 * time.Millisecond)
					err := hub.SendMessageWithResult(c.GetID(), []byte(strconv.Itoa(rn)))
					if err != nil {
						log.Println("SendMessageWithResult:", err)
						return
					}
				}
			}()
			num++
			if num > 200 {
				// hub.RemoveRandomClient()
				num--
			}
		}
	}()

	go func() {
		for {
			time.Sleep(33 * time.Millisecond)
			err := hub.Broadcast([]byte("hey"))
			if err != nil {
				log.Println("Broadcast:", err)
				return
			}
		}
	}()

	time.Sleep(3 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	hub.Wait()
	// time.Sleep(12 * time.Second)
	leftOver := hub.LenClients()
	log.Println("stopped")
	if leftOver != 0 {
		log.Println("left over:", leftOver)
		newLeftOver := hub.LenClients()
		log.Println("new left over:", newLeftOver)
	}

	return leftOver
}
func TestReconnects(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	for i := 0; i < 10; i++ {
		assert.EqualValues(t, 0, test())
	}
}

func testWithoutBroadcast() int {
	hub := siwebsocket.NewHub("http://127.0.0.1:8080",
		"/path/_push", 10*time.Second, 60*time.Second, 1024000, true,
		siwebsocket.WithAfterStoreClient(func(c sicore.Client, err error) {
			if err != nil {
				log.Println("store: "+err.Error(), c.GetID())
			} else {
				log.Println("store: ", c.GetID())
			}
		}),
		siwebsocket.WithAfterDeleteClient(func(c sicore.Client, err error) {
			if err != nil {
				log.Println("delete: "+err.Error(), c.GetID())
			} else {
				log.Println("delete: ", c.GetID())
			}
		}))
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push/randomclose"}

	go func() {
		num := 0
		for {
			time.Sleep(80 * time.Millisecond)
			log.Println("connect")
			conn, _, err := siwebsocket.DefaultConn(u, nil)
			if err != nil {
				log.Println(err)
				return
			}

			c, err := siwebsocket.NewClient(conn,
				siwebsocket.WithWriteWait(10*time.Second),
				siwebsocket.WithReadWait(60*time.Second),
				siwebsocket.WithMaxMessageSize(1024000),
				siwebsocket.WithUsePingPong(true),
				siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageLogHandler{}),
				siwebsocket.WithHub(hub),
				siwebsocket.WithUserID("9099909"),
				siwebsocket.WithUserGroupID("90999"),
			)
			if err != nil {
				log.Println(err)
				return
			}

			rn := rand.Intn(1000)
			if rn == 0 {
				hub.RemoveRandomClient()
			}
			go func() {
				for {
					time.Sleep(50 * time.Millisecond)
					err := hub.SendMessageWithResult(c.GetID(), []byte(strconv.Itoa(rn)))
					if err != nil {
						log.Println("SendMessageWithResult:", err)
						return
					}
				}
			}()
			num++
			if num > 200 {
				// hub.RemoveRandomClient()
				num--
			}
		}
	}()

	time.Sleep(4 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	hub.Wait()
	// time.Sleep(12 * time.Second)
	leftOver := hub.LenClients()
	log.Println("stopped", leftOver)

	return leftOver
}

func TestReconnectsWithoutBroadcast(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	for i := 0; i < 5; i++ {
		assert.EqualValues(t, 0, testWithoutBroadcast())
	}
}

func testBroadcast() int {
	hub := siwebsocket.NewHub("http://127.0.0.1:8080", "/path/_push", 10*time.Second, 60*time.Second, 1024000, true)
	go hub.Run()

	u := url.URL{Scheme: "ws", Host: ":48080", Path: "/push/randomclose"}

	log.Println("start")
	go func() {
		for {
			time.Sleep(44 * time.Millisecond)
			conn, _, err := siwebsocket.DefaultConn(u, nil)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = siwebsocket.NewClient(conn,
				siwebsocket.WithWriteWait(10*time.Second),
				siwebsocket.WithReadWait(60*time.Second),
				siwebsocket.WithMaxMessageSize(1024000),
				siwebsocket.WithUsePingPong(true),
				siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageLogHandler{}),
				siwebsocket.WithHub(hub),
			)
			if err != nil {
				log.Println(err)
				return
			}

			rn := rand.Intn(100)
			if rn == 0 {
				hub.RemoveRandomClient()
			}
		}
	}()

	go func() {
		for {
			time.Sleep(33 * time.Millisecond)
			err := hub.Broadcast([]byte("hey"))
			if err != nil {
				log.Println("Broadcast:", err)
				return
			}
		}
	}()

	time.Sleep(4 * time.Second)
	log.Println("stopping...")
	hub.Stop()
	hub.Wait()
	// time.Sleep(12 * time.Second)
	leftOver := hub.LenClients()
	log.Println("stopped")
	if leftOver != 0 {
		log.Println("clients left over:", leftOver)
		log.Println("new len:", hub.LenClients())
	}
	return leftOver
}

func TestBroadcast(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}
	for i := 0; i < 5; i++ {
		assert.EqualValues(t, 0, testBroadcast())
	}
	// log.Println("waiting...")
	// time.Sleep(12 * time.Second)
}
