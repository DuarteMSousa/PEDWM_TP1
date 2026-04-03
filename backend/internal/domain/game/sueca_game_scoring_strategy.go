package game

import (
	"backend/internal/domain/round"
)

type SuecaGameScoringStrategy struct {
}

func NewSuecaGameScoringStrategy() SuecaGameScoringStrategy {
	return SuecaGameScoringStrategy{}
}

func (s SuecaGameScoringStrategy) CalculateCurrentRoundGamePoints(round *round.Round) map[string]int {
	points := make(map[string]int)

	roundScore := round.GetScore()

	for teamID := range round.Teams {
		teamGamePoints := 0

		if roundScore[teamID] == 120 {
			teamGamePoints = 4
		} else if roundScore[teamID] > 90 {
			teamGamePoints = 2
		} else if roundScore[teamID] >= 60 {
			teamGamePoints = 1
		} else {
			teamGamePoints = 0
		}

		points[teamID] = teamGamePoints
	}

	return points
}

func (s SuecaGameScoringStrategy) HasGameEnded(game *Game) bool {
	for _, team := range game.Teams {
		if game.Score[team.ID] >= 5 {
			return true
		}
	}
	return false
}

func (s SuecaGameScoringStrategy) Winner(game *Game) string {
	var winner string
	var maxPoints int
	for teamID := range game.Teams {
		if game.Score[teamID] > maxPoints {
			maxPoints = game.Score[teamID]
			winner = teamID
		}
	}
	return winner
}
