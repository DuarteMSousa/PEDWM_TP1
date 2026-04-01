package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/round"
)

// GamePlayingState implementa GameState
type GamePlayingState struct {
	game *Game
}

func NewGamePlayingState(g *Game) *GamePlayingState {
	return &GamePlayingState{game: g}
}

func (s *GamePlayingState) Enter() {
	s.game.Status = IN_PROGRESS
	s.game.State.Update()
}

func (s *GamePlayingState) Update() {
	s.game.UpdateRoundState()

	if s.game.round.RuleStrategy.HasEnded(s.game.round) {

		teamScores := s.game.scoringStrategy.CalculateCurrentRoundGamePoints(s.game.round)

		for teamID, points := range teamScores {
			s.game.Score[teamID] += points
		}
		s.game.AddEvent(events.NewGameScoreUpdatedEvent(s.game.ID.String(), s.game.Score))

		if s.game.scoringStrategy.HasGameEnded(s.game) {
			s.game.State = NewGameFinishedState(s.game)
			s.game.State.Enter()
			return
		}

		s.game.round = round.NewRound(s.game.ID, s.game.Teams, s.game.botStrategy)
		s.game.round.State.Enter()
		s.game.UpdateRoundState()
	}

	if s.game.scoringStrategy.HasGameEnded(s.game) {
		s.game.State = NewGameFinishedState(s.game)
		s.game.State.Enter()
		return
	}
}
