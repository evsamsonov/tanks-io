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

	var game = tanks.NewGame()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("Failed to create socket io server", err)
	}

	mainRoom := "main"
	server.OnConnect("/", func(s socketio.Conn) error {
		s.Join(mainRoom)

		player := game.AddPlayer(s.ID())
		log.Printf("Player #%s connected %+v", s.ID(), player)

		err = broadcastCurrentPlayers(game, server, mainRoom)
		if err != nil {
			log.Printf("Failed to broadcast player on connect: %v", err)
			return err
		}

		return nil
	})
	server.OnEvent("/", "playerMovement", func(s socketio.Conn, msg string) {
		payload := struct {
			X int `json:"x"`
			Y int `json:"y"`
		}{}
		err := json.Unmarshal([]byte(msg), &payload)
		if err != nil {
			log.Printf("Failed to unmarshal player movement payload: %v", err)
			return
		}

		err = game.MovePlayer(s.ID(), payload.X, payload.Y)
		if err != nil {
			log.Printf("Failed to move player: %v, %v", err, payload)
			return
		}

		err = broadcastCurrentPlayers(game, server, mainRoom)
		if err != nil {
			log.Printf("Failed to broadcast player: %v, %v", err, payload)
			return
		}
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
		server.BroadcastToRoom("/", mainRoom, "disconnect", string(disconnectPayload))

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

func broadcastCurrentPlayers(game *tanks.Game, server *socketio.Server, mainRoom string) error {
	currentPlayersPayload, err := json.Marshal(game.Players())
	if err != nil {
		log.Printf("Failed to marshal players: %v", err)
		return err
	}

	server.BroadcastToRoom("/", mainRoom, "currentPlayers", string(currentPlayersPayload))
	return nil
}
