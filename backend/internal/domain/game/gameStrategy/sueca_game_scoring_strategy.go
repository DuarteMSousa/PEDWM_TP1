package game_strategy

import (
	"backend/internal/domain/round"
)

type SuecaGameScoringStrategy struct {
}

func NewSuecaGameScoringStrategy() *SuecaGameScoringStrategy {
	return &SuecaGameScoringStrategy{}
}

func (s *SuecaGameScoringStrategy) CalculateRoundPoints(round *round.Round) map[string]int {
	panic("SuecaGameScoringStrategy.CalculateRoundPoints is not implemented yet")
}
