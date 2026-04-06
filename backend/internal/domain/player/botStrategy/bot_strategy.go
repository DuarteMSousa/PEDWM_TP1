package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

// BotStrategyType identifies the type of bot strategy.
type BotStrategyType string

const (
	EASY BotStrategyType = "EASY"
	HARD BotStrategyType = "HARD"
)

// IBotStrategy defines the interface for automatic game strategies for bots.
type IBotStrategy interface {
	ChooseCard(hand hand.Hand, leadSuit card.Suit, cardStrengthProvider CardStrengthProvider) card.Card
	GetType() BotStrategyType
}
