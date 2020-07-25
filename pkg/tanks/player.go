package tanks

import (
	"errors"
	"math/rand"
	"sync"
)

type Direction int

const (
	UpDirection    Direction = 0
	RightDirection Direction = 1
	LeftDirection  Direction = -1
	DownDirection  Direction = -2
)

// 0 top, 1 right, -2 down, -1 left

var ErrNotFound = errors.New("player not found")

type Player struct {
	ID        string    `json:"id"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"direction"`
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
	list map[string]*Player
}

func NewPlayers() *Players {
	return &Players{
		list: make(map[string]*Player),
	}
}

func (p *Players) Add(player Player) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.list[player.ID] = &player
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
			playerStream <- *player
		}
	}()
	return playerStream
}

// fixme rename
func (p *Players) ChangePlayer(ID string, x, y int, direction Direction) error {
	p.mu.RLock()
	player, ok := p.list[ID]
	p.mu.RUnlock()
	if !ok {
		return ErrNotFound
	}

	p.mu.Lock()
	player.X = x
	player.Y = y
	player.Direction = direction
	p.mu.Unlock()
	return nil
}
