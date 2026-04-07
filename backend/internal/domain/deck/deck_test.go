package deck

import (
	"backend/internal/domain/card"
	"errors"
	"testing"
)

func makeCard(t *testing.T, suit card.Suit, rank card.Rank) card.Card {
	t.Helper()

	c, err := card.NewCard(suit, rank)
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}
	return c
}

func TestDeckFirstDrawRemainingAndEmpty(t *testing.T) {
	t.Parallel()

	c1 := makeCard(t, card.Hearts, card.A)
	c2 := makeCard(t, card.Spades, card.K)
	d := NewDeck([]card.Card{c1, c2})

	first, err := d.First()
	if err != nil {
		t.Fatalf("expected First() without error, got %v", err)
	}
	if first.ID != c1.ID {
		t.Fatalf("expected top card %q, got %q", c1.ID, first.ID)
	}
	if d.Remaining() != 2 {
		t.Fatalf("expected remaining 2 after First, got %d", d.Remaining())
	}

	drawn, err := d.Draw()
	if err != nil {
		t.Fatalf("expected Draw() without error, got %v", err)
	}
	if drawn.ID != c1.ID {
		t.Fatalf("expected drawn card %q, got %q", c1.ID, drawn.ID)
	}
	if d.Remaining() != 1 {
		t.Fatalf("expected remaining 1, got %d", d.Remaining())
	}

	_, _ = d.Draw()
	_, err = d.Draw()
	if !errors.Is(err, ErrDeckEmpty) {
		t.Fatalf("expected ErrDeckEmpty, got %v", err)
	}
}

func TestDeckResetAndIsEmpty(t *testing.T) {
	t.Parallel()

	d := NewDeck([]card.Card{
		makeCard(t, card.Clubs, card.J),
	})
	if d.IsEmpty() {
		t.Fatal("expected deck to be non-empty")
	}

	d.Reset()
	if !d.IsEmpty() {
		t.Fatal("expected deck to be empty after Reset")
	}
	if d.Remaining() != 0 {
		t.Fatalf("expected remaining 0, got %d", d.Remaining())
	}
}
