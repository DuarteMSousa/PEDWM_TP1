package bot_strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/hand"
	"testing"
)

type fakeCardStrengthProvider struct{}

func (fakeCardStrengthProvider) CardStrength(rank card.Rank) int {
	switch rank {
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

func newCard(t *testing.T, suit card.Suit, rank card.Rank) card.Card {
	t.Helper()

	c, err := card.NewCard(suit, rank)
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}
	return c
}

func TestEasyBotChooseCardWithLeadSuit(t *testing.T) {
	t.Parallel()

	h := hand.NewHand()
	h.AddCard(newCard(t, card.Hearts, card.Two))
	h.AddCard(newCard(t, card.Clubs, card.A))

	strategy := NewEasyBotStrategy()
	chosen := strategy.ChooseCard(*h, card.Clubs, fakeCardStrengthProvider{})
	if chosen.Suit != card.Clubs {
		t.Fatalf("expected CLUBS card, got %s", chosen.Suit)
	}
}

func TestEasyBotChooseFirstWhenNoLeadSuitInHand(t *testing.T) {
	t.Parallel()

	h := hand.NewHand()
	first := newCard(t, card.Hearts, card.Three)
	h.AddCard(first)
	h.AddCard(newCard(t, card.Clubs, card.A))

	strategy := NewEasyBotStrategy()
	chosen := strategy.ChooseCard(*h, card.Spades, fakeCardStrengthProvider{})
	if chosen.ID != first.ID {
		t.Fatalf("expected first card %q, got %q", first.ID, chosen.ID)
	}
}

func TestHardBotChooseStrongestFromLeadSuit(t *testing.T) {
	t.Parallel()

	h := hand.NewHand()
	h.AddCard(newCard(t, card.Hearts, card.Two))
	best := newCard(t, card.Hearts, card.A)
	h.AddCard(best)
	h.AddCard(newCard(t, card.Clubs, card.K))

	strategy := NewHardBotStrategy()
	chosen := strategy.ChooseCard(*h, card.Hearts, fakeCardStrengthProvider{})
	if chosen.ID != best.ID {
		t.Fatalf("expected strongest lead card %q, got %q", best.ID, chosen.ID)
	}
}

func TestHardBotChooseStrongestOverallWhenNoLeadSuit(t *testing.T) {
	t.Parallel()

	h := hand.NewHand()
	h.AddCard(newCard(t, card.Clubs, card.Three))
	best := newCard(t, card.Hearts, card.A)
	h.AddCard(best)
	h.AddCard(newCard(t, card.Diamonds, card.K))

	strategy := NewHardBotStrategy()
	chosen := strategy.ChooseCard(*h, card.Spades, fakeCardStrengthProvider{})
	if chosen.ID != best.ID {
		t.Fatalf("expected strongest overall card %q, got %q", best.ID, chosen.ID)
	}
}
