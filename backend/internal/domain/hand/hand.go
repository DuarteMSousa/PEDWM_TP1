package hand

import (
	"backend/internal/domain/card"
	"errors"
)

var (
	ErrCardNotInHand = errors.New("Card not in hand")
)

// Hand represents a player's hand of cards.
type Hand struct {
	Cards []card.Card
}

// NewHand creates an empty hand.
func NewHand() *Hand {
	return &Hand{Cards: []card.Card{}}
}

// AddCard adds a card to the hand.
func (h *Hand) AddCard(c card.Card) {
	h.Cards = append(h.Cards, c)
}

// GetCard searches for a card by ID without removing it.
func (h *Hand) GetCard(cardId string) (card.Card, error) {
	for _, card := range h.Cards {
		if card.ID == cardId {
			return card, nil
		}
	}
	return card.Card{}, ErrCardNotInHand
}

// RemoveCard removes and returns a card by ID.
func (h *Hand) RemoveCard(cardId string) (card.Card, error) {
	for i, card := range h.Cards {
		if card.ID == cardId {
			h.Cards = append(h.Cards[:i], h.Cards[i+1:]...)
			return card, nil
		}
	}
	return card.Card{}, ErrCardNotInHand
}

// IsEmpty indicates if the hand is empty.
func (h *Hand) IsEmpty() bool {
	return len(h.Cards) == 0
}

// HasSuit indicates if the hand contains at least one card of the specified suit.
func (h *Hand) HasSuit(suit card.Suit) bool {
	for _, card := range h.Cards {
		if card.Suit == suit {
			return true
		}
	}
	return false
}
