package game

import "backend/internal/domain/events"

// GamePlayingState implementa GameState
type GamePlayingState struct {
	game *Game
}

func NewGamePlayingState(g *Game) *GamePlayingState {
	return &GamePlayingState{game: g}
}

func (s *GamePlayingState) Enter() {
	event := events.NewGameStartedEvent(s.game.ID.String())
	s.game.AddEvent(event)
	s.game.State.Update()
}

func (s *GamePlayingState) Update() {
	if s.game.round.RuleStrategy.HasEnded(s.game.round) {

		teamScores := s.game.scoringStrategy.CalculateCurrentRoundGamePoints(s.game.round)

		for teamID, points := range teamScores {
			s.game.Teams[teamID].GameScore += points
		}

		if s.game.scoringStrategy.HasGameEnded(s.game) {
			s.game.State = NewGameFinishedState(s.game)
			s.game.State.Enter()
		}
	}
}
