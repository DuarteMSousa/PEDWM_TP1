package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/round"
	"backend/internal/domain/team"
)

// GameStartingState implements GameState
type GameStartingState struct {
	game *Game
}

// NewGameStartingState creates a new starting game state.
func NewGameStartingState(g *Game) *GameStartingState {
	return &GameStartingState{game: g}
}

// Enter initializes the game, creates the first round, and transitions to the playing state.
func (s *GameStartingState) Enter() {
	teams := make([]team.Team, len(s.game.Teams))
	i := 0
	for _, t := range s.game.Teams {
		teams[i] = *t
		i++
	}

	s.game.AddEvent(events.NewGameStartedEvent(s.game.ID.String(), teams))
	if s.game.Score == nil {
		s.game.Score = make(map[string]int)
	}
	for _, t := range s.game.Teams {
		s.game.Score[t.ID] = 0
	}

	s.game.round = round.NewRound(s.game.ID, s.game.Teams, s.game.botStrategy)

	s.game.round.State.Enter()

	s.game.State.Update()
}

// Update transitions the game to the playing state.
func (s *GameStartingState) Update() {
	s.game.State = NewGamePlayingState(s.game)
	s.game.State.Enter()
}
