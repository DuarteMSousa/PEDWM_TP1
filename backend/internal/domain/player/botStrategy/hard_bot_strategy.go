package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
)

// HardBotStrategy implements an advanced strategy: plays the strongest card
// of the leading suit, or the strongest card in hand if it doesn't have the suit.
type HardBotStrategy struct{}

// NewHardBotStrategy creates a new instance of the hard strategy.
func NewHardBotStrategy() *HardBotStrategy {
	return &HardBotStrategy{}
}

// ChooseCard chooses the card to play according to the hard strategy.
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

// GetType returns the type of the strategy (HARD).
func (h *HardBotStrategy) GetType() BotStrategyType {
	return HARD
}
