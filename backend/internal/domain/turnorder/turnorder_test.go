package turnorder

import (
	"backend/internal/domain/player"
	"errors"
	"testing"
)

func players4() []*player.Player {
	return []*player.Player{
		player.NewPlayer("p3", "P3", 3),
		player.NewPlayer("p1", "P1", 1),
		player.NewPlayer("p4", "P4", 4),
		player.NewPlayer("p2", "P2", 2),
	}
}

func TestNewTurnOrderRejectsInvalidSize(t *testing.T) {
	t.Parallel()

	_, err := NewTurnOrder("p1", players4()[:3])
	if !errors.Is(err, ErrTurnOrderInvalidSize) {
		t.Fatalf("expected ErrTurnOrderInvalidSize, got %v", err)
	}
}

func TestNewTurnOrderRejectsDuplicateIDs(t *testing.T) {
	t.Parallel()

	ps := []*player.Player{
		player.NewPlayer("p1", "P1", 1),
		player.NewPlayer("p1", "P1b", 2),
		player.NewPlayer("p3", "P3", 3),
		player.NewPlayer("p4", "P4", 4),
	}
	_, err := NewTurnOrder("p1", ps)
	if !errors.Is(err, ErrTurnOrderDuplicateIDs) {
		t.Fatalf("expected ErrTurnOrderDuplicateIDs, got %v", err)
	}
}

func TestNewTurnOrderStartsFromLeaderAndFollowsSequence(t *testing.T) {
	t.Parallel()

	order, err := NewTurnOrder("p3", players4())
	if err != nil {
		t.Fatalf("failed to create turn order: %v", err)
	}

	want := []string{"p3", "p4", "p1", "p2"}
	for i := 0; i < len(want); i++ {
		next, err := order.Next()
		if err != nil {
			t.Fatalf("unexpected error on Next at step %d: %v", i, err)
		}
		if next != want[i] {
			t.Fatalf("unexpected next player at step %d: got %q want %q", i, next, want[i])
		}
		if _, err := order.Dequeue(); err != nil {
			t.Fatalf("unexpected error on Dequeue at step %d: %v", i, err)
		}
	}
}

func TestTurnOrderRemoveAndContains(t *testing.T) {
	t.Parallel()

	order, err := NewTurnOrder("p1", players4())
	if err != nil {
		t.Fatalf("failed to create turn order: %v", err)
	}

	if !order.Contains("p2") {
		t.Fatal("expected order to contain p2")
	}

	if err := order.Remove("p2"); err != nil {
		t.Fatalf("expected remove to succeed, got %v", err)
	}

	if order.Contains("p2") {
		t.Fatal("expected order not to contain p2 after removal")
	}

	if err := order.Remove("missing"); !errors.Is(err, ErrTurnOrderPlayerAbsent) {
		t.Fatalf("expected ErrTurnOrderPlayerAbsent, got %v", err)
	}
}
