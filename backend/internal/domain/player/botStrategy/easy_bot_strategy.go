package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

// EasyBotStrategy implements a simple strategy: plays the first card
// of the leading suit, or the first card in hand if it doesn't have the suit.
type EasyBotStrategy struct{}

// NewEasyBotStrategy creates a new instance of the easy strategy.
func NewEasyBotStrategy() *EasyBotStrategy {
	return &EasyBotStrategy{}
}

// ChooseCard chooses the card to play according to the easy strategy.
func (e *EasyBotStrategy) ChooseCard(hand hand.Hand, leadSuit card.Suit, cardStrengthProvider CardStrengthProvider) card.Card {
	if len(hand.Cards) == 0 {
		return card.Card{}
	}

	for _, c := range hand.Cards {
		if c.Suit == leadSuit {
			return c
		}
	}

	return hand.Cards[0]
}

// GetType returns the type of the strategy (EASY).
func (e *EasyBotStrategy) GetType() BotStrategyType {
	return EASY
}
