package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/", home)
	http.HandleFunc("/state/", state)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/state/")
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
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

type State struct {
	Stamp   int64    `json:"t"`
	Players []Player `json:"players"`
}

type Player struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var currentState = State{
	Players: []Player{
		{X: 0.123, Y: 20},
		{X: 125.990, Y: 40},
	},
}

func state(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Failed to upgrade: ", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	for {
		currentState.Stamp = time.Now().UnixNano() / int64(time.Millisecond)
		err := conn.WriteJSON(currentState)
		if err != nil {
			log.Print("Failed to write json: ", err)
			break
		}
		<-ticker.C
	}
	ticker.Stop()
}
