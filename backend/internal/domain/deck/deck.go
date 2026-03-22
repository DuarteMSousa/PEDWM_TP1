package deck

import (
	"backend/internal/domain/card"
	"errors"
)

var (
	ErrDeckEmpty = errors.New("deck is empty")
)

type Deck struct {
	Cards []card.Card
}

func NewDeck() *Deck {
	return &Deck{Cards: []card.Card{}}
}

func (d *Deck) Shuffle() {
	// Implementar lógica de embaralhamento

}

func (d *Deck) Draw() (card.Card, error) {
	if len(d.Cards) == 0 {
		return card.Card{}, ErrDeckEmpty
	}
	topCard := d.Cards[0]
	d.Cards = d.Cards[1:]
	return topCard, nil
}

func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}

func (d *Deck) Reset() {
	d.Cards = []card.Card{}
}

func (d *Deck) Remaining(c card.Card) int {
	return len(d.Cards)
}
