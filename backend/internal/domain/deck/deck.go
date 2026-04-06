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

// Deck represents a deck of cards.
type Deck struct {
	cards []card.Card
}

// NewDeck creates a deck from a set of cards.
func NewDeck(cards []card.Card) *Deck {
	return &Deck{cards: cards}
}

// Shuffle shuffles the cards randomly.
func (d *Deck) Shuffle() {
	if len(d.cards) < 2 {
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})

}

// First returns the top card without removing it.
func (d *Deck) First() (card.Card, error) {
	if len(d.cards) == 0 {
		return card.Card{}, ErrDeckEmpty
	}
	topCard := d.cards[0]
	return topCard, nil
}

// Draw removes and returns the top card.
func (d *Deck) Draw() (card.Card, error) {
	if len(d.cards) == 0 {
		return card.Card{}, ErrDeckEmpty
	}
	topCard := d.cards[0]
	d.cards = d.cards[1:]
	return topCard, nil
}

// IsEmpty indicates whether the deck is empty.
func (d *Deck) IsEmpty() bool {
	return len(d.cards) == 0
}

// Reset clears the deck.
func (d *Deck) Reset() {
	d.cards = []card.Card{}
}

// Remaining returns the number of cards remaining in the deck.
func (d *Deck) Remaining() int {
	return len(d.cards)
}
