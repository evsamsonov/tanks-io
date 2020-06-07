package game

type State struct {
	Stamp   int64    `json:"t"`
	Players []Player `json:"players"`
}

type Player struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
