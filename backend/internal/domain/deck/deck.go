package deck

import (
	"backend/internal/domain/card"
	"errors"
	"math/rand"
	"time"
)

var (
	ErrDeckEmpty = errors.New("deck is empty")
)

type Deck struct {
	cards []card.Card
}

func NewDeck() *Deck {
	return &Deck{cards: []card.Card{}}
}

func (d *Deck) Shuffle() {
	if len(d.cards) < 2 {
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})

}

func (d *Deck) First() (card.Card, error) {
	if len(d.cards) == 0 {
		return card.Card{}, ErrDeckEmpty
	}
	topCard := d.cards[0]
	return topCard, nil
}

func (d *Deck) Draw() (card.Card, error) {
	if len(d.cards) == 0 {
		return card.Card{}, ErrDeckEmpty
	}
	topCard := d.cards[0]
	d.cards = d.cards[1:]
	return topCard, nil
}

func (d *Deck) IsEmpty() bool {
	return len(d.cards) == 0
}

func (d *Deck) Reset() {
	d.cards = []card.Card{}
}

func (d *Deck) Remaining(c card.Card) int {
	return len(d.cards)
}
