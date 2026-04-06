package game

import (
	"backend/internal/domain/round"
)

// IGameScoringStrategy defines the strategy for game-level scoring.
type IGameScoringStrategy interface {
	CalculateCurrentRoundGamePoints(round *round.Round) map[string]int
	HasGameEnded(game *Game) bool
	Winner(game *Game) string
}
