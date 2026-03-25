package game

import (
	"backend/internal/domain/round"
)

type SuecaGameScoringStrategy struct {
}

func NewSuecaGameScoringStrategy() SuecaGameScoringStrategy {
	return SuecaGameScoringStrategy{}
}

func (s *SuecaGameScoringStrategy) CalculateRoundGamePoints(round *round.Round) map[string]int {
	points := make(map[string]int)

	for teamID, team := range round.Teams {
		teamGamePoints := 0

		if team.RoundScore == 120 {
			teamGamePoints = 4
		} else if team.RoundScore >= 60 {
			teamGamePoints = 1
		} else {
			teamGamePoints = 0
		}

		points[teamID] = teamGamePoints
	}

	return points
}

func (s *SuecaGameScoringStrategy) HasGameEnded(game *Game) bool {
	for _, team := range game.Teams {
		if team.GameScore >= 4 {
			return true
		}
	}
	return false
}

func (s *SuecaGameScoringStrategy) Winner(game *Game) string {
	var winner string
	var maxPoints int
	for teamID, team := range game.Teams {
		if team.GameScore > maxPoints {
			maxPoints = team.GameScore
			winner = teamID
		}
	}
	return winner
}
