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
	bestStrength := trickStrength(best.Rank)
	for _, c := range candidates[1:] {
		s := trickStrength(c.Rank)
		if s > bestStrength {
			best = c
			bestStrength = s
		}
	}

	return best
}

func trickStrength(r card.Rank) int {
	switch r {
	case card.A:
		return 10
	case card.Seven:
		return 9
	case card.K:
		return 8
	case card.J:
		return 7
	case card.Q:
		return 6
	case card.Six:
		return 5
	case card.Five:
		return 4
	case card.Four:
		return 3
	case card.Three:
		return 2
	case card.Two:
		return 1
	default:
		return 0
	}
}
