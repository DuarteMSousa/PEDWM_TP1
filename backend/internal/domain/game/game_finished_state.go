package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/team"
)

// GameFinishedState implements GameState
type GameFinishedState struct {
	game *Game
}

// NewGameFinishedState creates a new finished game state.
func NewGameFinishedState(game *Game) *GameFinishedState {
	return &GameFinishedState{game: game}
}

// Enter sets the game status to FINISHED and updates the state.
func (s *GameFinishedState) Enter() {
	s.game.Status = FINISHED
	s.game.State.Update()
}

// Update calculates the winner and adds a GameEndedEvent to the game's events.
func (s *GameFinishedState) Update() {
	winner := s.game.scoringStrategy.Winner(s.game)
	teams := make([]team.Team, len(s.game.Teams))
	i := 0
	for _, t := range s.game.Teams {
		teams[i] = *t
		i++
	}

	event := events.NewGameEndedEvent(s.game.ID.String(), s.game.Score, winner, teams)
	s.game.AddEvent(event)
}
