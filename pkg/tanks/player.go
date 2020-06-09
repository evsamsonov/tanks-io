package tanks

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Player struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func NewPlayerWithRandLoc(ID string) Player {
	return Player{
		ID: ID,
		X:  rand.Intn(750),
		Y:  rand.Intn(550),
	}
}

type Players struct {
	mu   sync.RWMutex
	list map[string]Player
}

func NewPlayers() *Players {
	return &Players{
		list: make(map[string]Player),
	}
}

func (p *Players) Add(player Player) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.list[player.ID] = player
}

func (p *Players) Remove(ID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.list, ID)
}

func (p *Players) All() <-chan Player {
	p.mu.RLock()
	defer p.mu.RUnlock()

	playerStream := make(chan Player)
	go func() {
		defer close(playerStream)
		for _, player := range p.list {
			playerStream <- player
		}
	}()
	return playerStream
}
