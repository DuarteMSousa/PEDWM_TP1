package hand

import (
	"backend/internal/domain/card"
	"errors"
	"testing"
)

func mustCard(t *testing.T, suit card.Suit, rank card.Rank) card.Card {
	t.Helper()

	c, err := card.NewCard(suit, rank)
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}
	return c
}

func TestHandAddGetRemoveLifecycle(t *testing.T) {
	t.Parallel()

	h := NewHand()
	c := mustCard(t, card.Hearts, card.A)

	h.AddCard(c)
	if h.IsEmpty() {
		t.Fatal("expected hand not to be empty after AddCard")
	}

	found, err := h.GetCard(c.ID)
	if err != nil {
		t.Fatalf("expected to find card, got error: %v", err)
	}
	if found.ID != c.ID {
		t.Fatalf("unexpected card found: got %q want %q", found.ID, c.ID)
	}

	removed, err := h.RemoveCard(c.ID)
	if err != nil {
		t.Fatalf("expected to remove card, got error: %v", err)
	}
	if removed.ID != c.ID {
		t.Fatalf("unexpected removed card: got %q want %q", removed.ID, c.ID)
	}
	if !h.IsEmpty() {
		t.Fatal("expected hand to be empty after removing the only card")
	}
}

func TestHandRemoveCardMissing(t *testing.T) {
	t.Parallel()

	h := NewHand()
	_, err := h.RemoveCard("missing")
	if !errors.Is(err, ErrCardNotInHand) {
		t.Fatalf("expected ErrCardNotInHand, got %v", err)
	}
}

func TestHandHasSuit(t *testing.T) {
	t.Parallel()

	h := NewHand()
	h.AddCard(mustCard(t, card.Spades, card.K))
	h.AddCard(mustCard(t, card.Clubs, card.Two))

	if !h.HasSuit(card.Spades) {
		t.Fatal("expected hand to have SPADES")
	}
	if h.HasSuit(card.Hearts) {
		t.Fatal("expected hand not to have HEARTS")
	}
}
