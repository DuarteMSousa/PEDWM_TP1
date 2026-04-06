package game

import (
	"backend/internal/domain/round"
)

// SuecaGameScoringStrategy implements the official scoring rules of Sueca:
// 120 points = 4 game points, >90 = 2 game points, >=60 = 1 game point, <60 = 0 game points.
// The game ends when a team reaches 4 game points.
type SuecaGameScoringStrategy struct {
}

// NewSuecaGameScoringStrategy creates a new instance of the Sueca scoring strategy.
func NewSuecaGameScoringStrategy() SuecaGameScoringStrategy {
	return SuecaGameScoringStrategy{}
}

// CalculateCurrentRoundGamePoints calculates the game points to be awarded
// to each team based on the round's score.
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

// HasGameEnded checks if any team has reached 4 game points.
func (s SuecaGameScoringStrategy) HasGameEnded(game *Game) bool {
	for _, team := range game.Teams {
		if game.Score[team.ID] >= 4 {
			return true
		}
	}
	return false
}

// Winner returns the ID of the team with the most game points.
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
