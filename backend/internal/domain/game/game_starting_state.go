package game

import (
	"backend/internal/domain/events"
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
	s.game.AddEvent(events.NewGameStartedEvent(s.game.ID.String()))
	for _, t := range s.game.Teams {
		s.game.Score[t.ID] = 0
	}

	s.game.round = round.NewRound(s.game.ID, s.game.Teams, s.game.botStrategy)

	s.game.round.State.Enter()

	s.game.State.Update()
}

func (s *GameStartingState) Update() {
	s.game.State = NewGamePlayingState(s.game)
	s.game.State.Enter()
}
