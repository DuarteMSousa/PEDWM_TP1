package hand

import (
	"backend/internal/domain/card"
	"errors"
)

var (
	ErrCardNotInHand = errors.New("Card not in hand")
)

type Hand struct {
	Cards []card.Card
}

func NewHand() *Hand {
	return &Hand{Cards: []card.Card{}}
}

func (h *Hand) AddCard(c card.Card) {
	h.Cards = append(h.Cards, c)
}

func (h *Hand) GetCard(cardId string) (card.Card, error) {
	for _, card := range h.Cards {
		if card.ID == cardId {
			return card, nil
		}
	}
	return card.Card{}, ErrCardNotInHand
}

func (h *Hand) RemoveCard(cardId string) (card.Card, error) {
	for i, card := range h.Cards {
		if card.ID == cardId {
			h.Cards = append(h.Cards[:i], h.Cards[i+1:]...)
			return card, nil
		}
	}
	return card.Card{}, ErrCardNotInHand
}

func (h *Hand) IsEmpty() bool {
	return len(h.Cards) == 0
}

func (h *Hand) HasSuit(suit card.Suit) bool {
	for _, card := range h.Cards {
		if card.Suit == suit {
			return true
		}
	}
	return false
}
