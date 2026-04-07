package card

import (
	"errors"
	"testing"
)

func TestNewCardValid(t *testing.T) {
	t.Parallel()

	c, err := NewCard(Hearts, A)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c.ID != "A_HEARTS" {
		t.Fatalf("unexpected card ID: got %q", c.ID)
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("expected card to validate, got %v", err)
	}
}

func TestNewCardInvalidSuit(t *testing.T) {
	t.Parallel()

	_, err := NewCard(Suit("INVALID"), A)
	if !errors.Is(err, ErrInvalidSuit) {
		t.Fatalf("expected ErrInvalidSuit, got %v", err)
	}
}

func TestNewCardInvalidRank(t *testing.T) {
	t.Parallel()

	_, err := NewCard(Hearts, Rank("X"))
	if !errors.Is(err, ErrInvalidRank) {
		t.Fatalf("expected ErrInvalidRank, got %v", err)
	}
}

func TestCardValidateInvalidID(t *testing.T) {
	t.Parallel()

	c := Card{
		ID:   " ",
		Suit: Hearts,
		Rank: A,
	}

	err := c.Validate()
	if !errors.Is(err, ErrInvalidCardID) {
		t.Fatalf("expected ErrInvalidCardID, got %v", err)
	}
}

func TestCardIsTrump(t *testing.T) {
	t.Parallel()

	c, err := NewCard(Spades, K)
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}

	if !c.IsTrump(Spades) {
		t.Fatal("expected card to be trump")
	}
	if c.IsTrump(Hearts) {
		t.Fatal("expected card not to be trump for HEARTS")
	}
}
