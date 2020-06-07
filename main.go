package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("Failed to create socket io server", err)
	}
	server.OnConnect("/", func(conn socketio.Conn) error {
		//conn.SetContext()
		fmt.Println("connected", conn.ID())
		return nil
	})
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal("Failed to serve socket io", err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			log.Fatal("Failed to close socket io ", err)
		}
	}()

	http.Handle("/state/", server)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	log.Println("Serving at " + *addr + "...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type State struct {
	Stamp   int64    `json:"t"`
	Players []Player `json:"players"`
}

type Player struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
