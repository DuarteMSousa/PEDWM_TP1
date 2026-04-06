package bot_strategy

import "backend/internal/domain/card"

// CardStrengthProvider abstracts the calculation of a card's strength.
type CardStrengthProvider interface {
	CardStrength(card.Rank) int
}
