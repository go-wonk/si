package main

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-wonk/si/siwebsocket"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":28080", "http service address")
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// func echo(w http.ResponseWriter, r *http.Request) {
// 	userID := r.URL.Query().Get("user_id")
// 	userGroupID := r.URL.Query().Get("user_group_id")

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade:", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}
// 	c, err := hub.CreateAndAddClient(conn,
// 		siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageLogHandler{}),
// 	)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}
// 	log.Println(c.GetID(), userID, userGroupID)
// }

func echoHub(hub *siwebsocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		userGroupID := r.URL.Query().Get("user_group_id")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		c, err := hub.CreateAndAddClient(conn,
			siwebsocket.WithMessageHandler(&siwebsocket.DefaultMessageLogHandler{}),
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Println(c.GetID(), userID, userGroupID)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := siwebsocket.NewHub("http://127.0.0.1:8080", "/ws/_push", 10*time.Second, 60*time.Second, 1024000, true)
	go hub.Run()
	defer func() {
		hub.Stop()
		hub.Wait()
	}()

	http.Handle("/ws/echo", echoHub(hub))
	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(*addr, nil))
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
