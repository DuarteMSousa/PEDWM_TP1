package game

import (
	"fmt"
)

// GameFinishedState implementa GameState
type GameFinishedState struct{}

func (s *GameFinishedState) Enter()  { fmt.Println("Fim de jogo.") }
func (s *GameFinishedState) Update() { /* Lógica de placar final */ }
