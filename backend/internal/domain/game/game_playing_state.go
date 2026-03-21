package game

// GamePlayingState implementa GameState
type GamePlayingState struct {
	game *Game
}

func NewGamePlayingState(g *Game) *GamePlayingState {
	return &GamePlayingState{game: g}
}

func (s *GamePlayingState) Enter() {

}
func (s *GamePlayingState) Update() {

}
