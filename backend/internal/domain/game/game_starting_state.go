package game

import (
	"backend/internal/domain/round"
)

// GameStartingState implementa GameState
type GameStartingState struct {
	game *Game
}

func NewGameStartingState(g *Game) *GameStartingState {
	return &GameStartingState{game: g}
}

func (s *GameStartingState) Enter() {
	for _, t := range s.game.teams {
		t.GameScore = 0
	}

	s.game.round = round.NewRound(s.game.teams)

	s.game.round.State.Enter()
}

func (s *GameStartingState) Update() {

}
