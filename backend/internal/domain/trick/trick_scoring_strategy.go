package trick

import "backend/internal/domain/card"

type ITrickScoringStrategy interface {
	CardPoints(card card.Card) int
	TrickPoints(plays []Play) int
}
