package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

type BotStrategyType string

const (
	EASY BotStrategyType = "EASY"
	HARD BotStrategyType = "HARD"
)

type IBotStrategy interface {
	ChooseCard(hand hand.Hand, leadSuit card.Suit, cardStrengthProvider CardStrengthProvider) card.Card
	GetType() BotStrategyType
}
