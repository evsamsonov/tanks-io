package tanks

import (
	"math/rand"
	"time"
)

type Game struct {
	players *Players
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	return &Game{
		players: NewPlayers(),
	}
}

func (g *Game) AddPlayer(ID string) Player {
	player := NewPlayerWithRandLoc(ID)
	g.players.Add(player)
	return player
}

func (g *Game) Players() []Player {
	players := make([]Player, 0)
	for player := range g.players.All() {
		players = append(players, player)
	}
	return players
}

func (g *Game) RemovePlayer(ID string) {
	g.players.Remove(ID)
}

func (g *Game) MovePlayer(ID string, x, y int, direction Direction) error {
	err := g.players.ChangePlayer(ID, x, y, direction)
	if err != nil {
		return err
	}
	return nil
}
