package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

type EasyBotStrategy struct{}

func NewEasyBotStrategy() *EasyBotStrategy {
	return &EasyBotStrategy{}
}

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
