package card

import (
	"errors"
	"fmt"
	"strings"
)

// Suit represents the suit of a card.
type Suit string

// Rank represents the face value of a card.
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

// Suits contains all valid suits.
var Suits = []Suit{Hearts, Spades, Diamonds, Clubs}

// Ranks contains all valid ranks, ordered by descending strength in Sueca.
var Ranks = []Rank{A, K, Q, J, Seven, Six, Five, Four, Three, Two}

var (
	ErrInvalidSuit   = errors.New("invalid suit")
	ErrInvalidRank   = errors.New("invalid rank")
	ErrInvalidCardID = errors.New("invalid card id")
)

// Valid checks if the suit is one of the allowed values.
func (n Suit) Valid() bool {
	switch n {
	case Hearts, Spades, Diamonds, Clubs:
		return true
	default:
		return false
	}
}

// Valid checks if the rank is one of the allowed values.
func (r Rank) Valid() bool {
	switch r {
	case A, K, Q, J, Seven, Six, Five, Four, Three, Two:
		return true
	default:
		return false
	}
}

// Card represents a playing card with an ID, suit, and rank.
type Card struct {
	ID   string `json:"id"`
	Suit Suit   `json:"suit"`
	Rank Rank   `json:"rank"`
}

// NewCard creates a new card, validating the suit and rank. The ID is generated
// automatically in the format "RANK_SUIT" (e.g., "A_HEARTS").
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

// Validate checks if the card has a valid ID, suit, and rank.
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

// IsTrump checks if the card belongs to the indicated trump suit.
func (c Card) IsTrump(trump Suit) bool {
	return c.Suit == trump
}
