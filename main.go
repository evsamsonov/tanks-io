package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"bitbucket.org/tanks-io/pkg/tanks"
	socketio "github.com/googollee/go-socket.io"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	game := tanks.NewGame()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("Failed to create socket io server", err)
	}

	mainRoom := "main"
	server.OnConnect("/", func(s socketio.Conn) error {
		s.Join(mainRoom)

		player := game.AddPlayer(s.ID())
		log.Printf("Player #%s connected %+v", s.ID(), player)

		currentPlayersPayload, err := json.Marshal(game.Players())
		if err != nil {
			log.Printf("Failed to marshal players: %v", err)
			return err
		}

		newPlayerPayload, err := json.Marshal(player)
		if err != nil {
			log.Printf("Failed to marshal players: %v", err)
			return err
		}

		s.Emit("currentPlayers", string(currentPlayersPayload))
		server.BroadcastToRoom("", mainRoom, "newPlayer", string(newPlayerPayload))

		return nil
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		game.RemovePlayer(s.ID())

		disconnect := struct {
			ID string `json:"id"`
		}{ID: s.ID()}
		disconnectPayload, err := json.Marshal(disconnect)
		if err != nil {
			log.Printf("Failed to marshal players: %v", err)
			return
		}
		server.BroadcastToRoom("", mainRoom, "disconnect", string(disconnectPayload))

		log.Printf("Payer #%s disconnected. Reason: %s", s.ID(), reason)
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
