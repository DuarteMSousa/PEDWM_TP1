package game

import (
	"backend/internal/domain/events"
)

// GameFinishedState implementa GameState
type GameFinishedState struct {
	game *Game
}

func NewGameFinishedState(game *Game) *GameFinishedState {
	return &GameFinishedState{game: game}
}

func (s *GameFinishedState) Enter() {
	s.game.Status = FINISHED
	s.game.State.Update()
}

func (s *GameFinishedState) Update() {
	winner := s.game.scoringStrategy.Winner(s.game)
	event := events.NewGameEndedEvent(s.game.ID.String(), s.game.Score, winner)
	s.game.AddEvent(event)
}
