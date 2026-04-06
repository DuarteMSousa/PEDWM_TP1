package trick

import "backend/internal/domain/card"

// ITrickScoringStrategy defines the scoring of cards and tricks.
type ITrickScoringStrategy interface {
	CardPoints(card card.Card) int
	TrickPoints(plays []Play) int
}
