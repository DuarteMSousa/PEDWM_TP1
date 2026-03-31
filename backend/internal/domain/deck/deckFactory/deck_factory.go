package deckFactory

import (
	"backend/internal/domain/card"
	"backend/internal/domain/deck"
)

func CreateSuecaDeck() *deck.Deck {

	cards := make([]card.Card, 0, 40)
	for _, suit := range card.Suits {
		for _, rank := range card.Ranks {
			createdCard, err := card.NewCard(suit, rank)
			if err != nil {
				panic(err)
			}
			cards = append(cards, createdCard)
		}
	}

	deck := deck.NewDeck(cards)

	return deck
}
