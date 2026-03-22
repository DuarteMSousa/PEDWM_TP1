package game_strategy

import (
	"backend/internal/domain/round"
)

type IGameScoringStrategy interface {
	CalculateRoundPoints(round *round.Round) map[string]int
}
