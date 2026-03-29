package game

// GamePlayingState implementa GameState
type GamePlayingState struct {
	game *Game
}

func NewGamePlayingState(g *Game) *GamePlayingState {
	return &GamePlayingState{game: g}
}

func (s *GamePlayingState) Enter() {
	s.game.State.Update()
}

func (s *GamePlayingState) Update() {
	s.game.UpdateRoundState()

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
