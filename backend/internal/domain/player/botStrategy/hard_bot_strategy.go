package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

type HardBotStrategy struct{}

func NewHardBotStrategy() *HardBotStrategy {
	return &HardBotStrategy{}
}

func (h *HardBotStrategy) ChooseCard(hand hand.Hand, leadSuit card.Suit) card.Card {

	panic("chooose card not implemented yet")
}
