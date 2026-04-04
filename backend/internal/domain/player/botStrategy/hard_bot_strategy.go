package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

type HardBotStrategy struct{}

func NewHardBotStrategy() *HardBotStrategy {
	return &HardBotStrategy{}
}

func (h *HardBotStrategy) ChooseCard(hand hand.Hand, leadSuit card.Suit, cardStrengthProvider CardStrengthProvider) card.Card {
	if len(hand.Cards) == 0 {
		return card.Card{}
	}

	candidates := make([]card.Card, 0, len(hand.Cards))
	for _, c := range hand.Cards {
		if c.Suit == leadSuit {
			candidates = append(candidates, c)
		}
	}
	if len(candidates) == 0 {
		candidates = hand.Cards
	}

	best := candidates[0]
	bestStrength := cardStrengthProvider.CardStrength(best.Rank)
	for _, c := range candidates[1:] {
		s := cardStrengthProvider.CardStrength(c.Rank)
		if s > bestStrength {
			best = c
			bestStrength = s
		}
	}

	return best
}

func (h *HardBotStrategy) GetType() BotStrategyType {
	return HARD
}
