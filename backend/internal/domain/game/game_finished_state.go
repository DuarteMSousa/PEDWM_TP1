package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/team"
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
	teams := make(map[string]team.Team)
	for _, team := range s.game.Teams {
		teams[team.ID] = *team
	}

	event := events.NewGameEndedEvent(s.game.ID.String(), s.game.Score, winner, teams)
	s.game.AddEvent(event)
}
