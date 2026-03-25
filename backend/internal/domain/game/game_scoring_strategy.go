package game

import (
	"backend/internal/domain/round"
)

type IGameScoringStrategy interface {
	CalculateCurrentRoundGamePoints(round *round.Round) map[string]int
	HasGameEnded(game *Game) bool
	Winner(game *Game) string
}
