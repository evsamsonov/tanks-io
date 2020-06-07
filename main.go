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
		conn.SetContext("")
		fmt.Println("connected", conn.ID())
		conn.Join("main")
		return nil
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", s.ID(), reason)
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

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	log.Println("Serving at " + *addr + "...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
