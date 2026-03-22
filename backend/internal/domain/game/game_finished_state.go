package game

import (
	"fmt"
)

// GameFinishedState implementa GameState
type GameFinishedState struct {
	game *Game
}

func NewGameFinishedState(game *Game) *GameFinishedState {
	return &GameFinishedState{game: game}
}

func (s *GameFinishedState) Enter()  { fmt.Println("Fim de jogo.") }
func (s *GameFinishedState) Update() { /* Lógica de placar final */ }
