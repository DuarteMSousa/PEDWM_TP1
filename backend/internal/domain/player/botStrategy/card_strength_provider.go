package bot_strategy

import "backend/internal/domain/card"

type CardStrengthProvider interface {
	CardStrength(card.Rank) int
}
