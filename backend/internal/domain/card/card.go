package card

import (
	"errors"
	"fmt"
	"strings"
)

type Suit string
type Rank string

const (
	Hearts   Suit = "HEARTS"
	Spades   Suit = "SPADES"
	Diamonds Suit = "DIAMONDS"
	Clubs    Suit = "CLUBS"
)

const (
	A     Rank = "A"
	K     Rank = "K"
	Q     Rank = "Q"
	J     Rank = "J"
	Seven Rank = "7"
	Six   Rank = "6"
	Five  Rank = "5"
	Four  Rank = "4"
	Three Rank = "3"
	Two   Rank = "2"
)

var Suits = []Suit{Hearts, Spades, Diamonds, Clubs}
var Ranks = []Rank{A, K, Q, J, Seven, Six, Five, Four, Three, Two}

var (
	ErrInvalidSuit   = errors.New("invalid suit")
	ErrInvalidRank   = errors.New("invalid rank")
	ErrInvalidCardID = errors.New("invalid card id")
)

func (n Suit) Valid() bool {
	switch n {
	case Hearts, Spades, Diamonds, Clubs:
		return true
	default:
		return false
	}
}

func (r Rank) Valid() bool {
	switch r {
	case A, K, Q, J, Seven, Six, Five, Four, Three, Two:
		return true
	default:
		return false
	}
}

// TrickStrength returns the Sueca trick order strength.
// Higher value means stronger card for deciding trick winner.
func (r Rank) TrickStrength() int {
	switch r {
	case A:
		return 10
	case Seven:
		return 9
	case K:
		return 8
	case J:
		return 7
	case Q:
		return 6
	case Six:
		return 5
	case Five:
		return 4
	case Four:
		return 3
	case Three:
		return 2
	case Two:
		return 1
	default:
		return 0
	}
}

type Card struct {
	ID   string `json:"id"`
	Suit Suit   `json:"suit"`
	Rank Rank   `json:"rank"`
}

func NewCard(Suit Suit, rank Rank) (Card, error) {
	if !Suit.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidSuit, Suit)
	}
	if !rank.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidRank, rank)
	}

	id := string(rank) + "_" + string(Suit)

	return Card{ID: id, Suit: Suit, Rank: rank}, nil
}

func (c Card) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return ErrInvalidCardID
	}
	if !c.Suit.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidSuit, c.Suit)
	}
	if !c.Rank.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidRank, c.Rank)
	}
	return nil
}

func (c Card) IsTrump(trump Suit) bool {
	return c.Suit == trump
}
