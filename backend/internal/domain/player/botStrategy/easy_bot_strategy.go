package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

type EasyBotStrategy struct{}

func NewEasyBotStrategy() *EasyBotStrategy {
	return &EasyBotStrategy{}
}

func (e *EasyBotStrategy) ChooseCard(hand hand.Hand, leadSuit card.Suit) card.Card {

	panic("chooose card not implemented yet")
}
