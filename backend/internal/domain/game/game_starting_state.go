package game

// GameStartingState implementa GameState
type GameStartingState struct {
	game *Game
}

func NewGameStartingState(g *Game) *GameStartingState {
	return &GameStartingState{game: g}
}

func (s *GameStartingState) Enter() {

}

func (s *GameStartingState) Update() {

}
