package card

import (
	"errors"
	"fmt"
	"strings"
)

type Naipe string

type Rank string

const (
	Hearts   Naipe = "HEARTS"
	Spades   Naipe = "SPADES"
	Diamonds Naipe = "DIAMONDS"
	Clubs    Naipe = "CLUBS"

	Copas   Naipe = Hearts
	Espadas Naipe = Spades
	Ouros   Naipe = Diamonds
	Paus    Naipe = Clubs
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

var (
	ErrInvalidNaipe  = errors.New("invalid suit")
	ErrInvalidRank   = errors.New("invalid rank")
	ErrInvalidCardID = errors.New("invalid card id")
)

func (n Naipe) Valid() bool {
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

type Card struct {
	ID    string `json:"id"`
	Naipe Naipe  `json:"naipe"`
	Rank  Rank   `json:"rank"`
}

func NewCard(id string, naipe Naipe, rank Rank) (Card, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Card{}, ErrInvalidCardID
	}
	if !naipe.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidNaipe, naipe)
	}
	if !rank.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidRank, rank)
	}
	return Card{ID: id, Naipe: naipe, Rank: rank}, nil
}

func (c Card) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return ErrInvalidCardID
	}
	if !c.Naipe.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidNaipe, c.Naipe)
	}
	if !c.Rank.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidRank, c.Rank)
	}
	return nil
}

func (c Card) IsTrump(trump Naipe) bool {
	return c.Naipe == trump
}
