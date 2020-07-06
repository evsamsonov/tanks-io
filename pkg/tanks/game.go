package tanks

type Game struct {
	players *Players
}

func NewGame() *Game {
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

func (g *Game) MovePlayer(ID string, x, y int) error {
	err := g.players.SetCoordinate(ID, x, y)
	if err != nil {
		return err
	}
	return nil
}
