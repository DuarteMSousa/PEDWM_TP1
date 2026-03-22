package round

import (
	"fmt"
)

// RoundFinishedState implementa RoundState
type RoundFinishedState struct {
	round *Round
}

func NewRoundFinishedState(round *Round) *RoundFinishedState {
	return &RoundFinishedState{round: round}
}

func (s *RoundFinishedState) Enter()  { fmt.Println("Fim de rodada.") }
func (s *RoundFinishedState) Update() { /* Lógica de placar final */ }
