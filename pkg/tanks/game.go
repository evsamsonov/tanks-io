package tanks

type Game struct {
	players *Players
}

func NewGame() *Game {
	return &Game{
		players: NewPlayers(),
	}
}

func (g *Game) AddPlayer(ID string) {
	g.players.Add(NewPlayerWithRandLoc(ID))
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
