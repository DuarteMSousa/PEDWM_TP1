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
	s.game.State.Update()
}

func (s *GameFinishedState) Update() {
	winner := s.game.scoringStrategy.Winner(s.game)
	finalScores := make(map[string]int)
	for teamID, team := range s.game.Teams {
		finalScores[teamID] = team.GameScore
	}
	event := events.NewGameEndedEvent(s.game.ID.String(), finalScores, winner)
	s.game.AddEvent(event)
}
