//go:build ignore
// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":48080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func idle(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	cnt := 0
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// n := rand.Intn(1000)
		// if n == 0 {
		// 	log.Printf("recv: %s", message)
		// }
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
		cnt++
	}
	log.Println("num sent:", cnt)
}

func echo(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	cnt := 0
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// n := rand.Intn(1000)
		// if n == 0 {
		// 	log.Printf("recv: %s", message)
		// }
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
		cnt++
	}
	log.Println("num sent:", cnt)
}

func echoRandomClose(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	cnt := 0
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		n := rand.Intn(1000)
		if n == 0 {
			log.Printf("recv: %s", message)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
		cnt++
	}
	log.Println("num sent:", cnt)
}
func push(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	cnt := 0
	go func() {
		for {
			time.Sleep(77 * time.Millisecond)
			n := rand.Intn(10000)
			err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(n)))
			if err != nil {
				log.Println("write1:", err)
				break
			}

			cnt++

			// if cnt > 100 {
			// 	c.Close()
			// 	return
			// }
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Println(string(message))

	}
	log.Println("num sent:", cnt)
}

func pushStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	cnt := 0
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			n := rand.Intn(1000)
			pushData := `{"id":%d,"name":"wonk","email_address":"wonk@wonk.org"}`
			pushData = fmt.Sprintf(pushData, n)
			err := c.WriteMessage(websocket.TextMessage, []byte(pushData))
			if err != nil {
				log.Println("write1:", err)
				break
			}

			cnt++

			// if cnt > 100 {
			// 	c.Close()
			// 	return
			// }
		}
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		n := rand.Intn(1000)
		if string(message) == strconv.Itoa(n) {
			log.Printf("recv: %s", message)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write2:", err)
			break
		}
	}
	log.Println("num sent:", cnt)
}

var (
	connections int64
)

func pushRandomClose(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
		v := atomic.AddInt64(&connections, -1)
		log.Println("num connections:", v)
	}()
	atomic.AddInt64(&connections, 1)
	cnt := 0
	go func() {
		for {
			time.Sleep(77 * time.Millisecond)
			n := rand.Intn(1000)
			err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(n)))
			if err != nil {
				log.Println("write1:", err)
				break
			}

			cnt++

			if n == 0 {
				c.Close()
				return
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Println(string(message))

	}
	log.Println("num sent:", cnt)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/idle", idle)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/echo/randomclose", echoRandomClose)
	http.HandleFunc("/push", push)
	http.HandleFunc("/push/student", pushStudent)
	http.HandleFunc("/push/randomclose", pushRandomClose)
	http.HandleFunc("/", home)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{Addr: *addr}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case sig := <-sigs:
				switch sig {
				case syscall.SIGINT:
					fallthrough
				case syscall.SIGTERM:
					fallthrough
				default:
					log.Printf("signal %v", sig)
				}
				server.Shutdown(context.Background())
				return
			case t := <-ticker.C:
				log.Printf("Connections: %v, NoOfGR:%v, %v", connections, runtime.NumGoroutine(), t)
			}
		}
	}()

	log.Fatal(server.ListenAndServe())
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
