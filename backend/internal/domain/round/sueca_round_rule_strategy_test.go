package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"testing"

	"github.com/google/uuid"
)

func mkCard(t *testing.T, suit card.Suit, rank card.Rank) card.Card {
	t.Helper()

	c, err := card.NewCard(suit, rank)
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}
	return c
}

func makeRoundTeams() map[string]*team.Team {
	p1 := player.NewPlayer("p1", "P1", 1)
	p2 := player.NewPlayer("p2", "P2", 2)
	p3 := player.NewPlayer("p3", "P3", 3)
	p4 := player.NewPlayer("p4", "P4", 4)

	return map[string]*team.Team{
		"t1": {ID: "t1", Players: []*player.Player{p1, p3}},
		"t2": {ID: "t2", Players: []*player.Player{p2, p4}},
	}
}

func TestSuecaRoundRuleStrategyHasEnded(t *testing.T) {
	t.Parallel()

	r := NewRound(uuid.New(), makeRoundTeams(), bot_strategy.NewEasyBotStrategy())
	strategy := NewSuecaRoundRuleStrategy()

	if !strategy.HasEnded(r) {
		t.Fatal("expected round to be ended when all hands are empty")
	}

	teamP1 := r.Teams["t1"].Players[0]
	teamP1.Hand.AddCard(mkCard(t, card.Hearts, card.A))
	if strategy.HasEnded(r) {
		t.Fatal("expected round not to be ended when at least one hand has cards")
	}
}

func TestSuecaRoundRuleStrategyWinner(t *testing.T) {
	t.Parallel()

	r := NewRound(uuid.New(), makeRoundTeams(), bot_strategy.NewEasyBotStrategy())
	strategy := NewSuecaRoundRuleStrategy()

	// Not ended: should return the error string.
	r.Teams["t1"].Players[0].Hand.AddCard(mkCard(t, card.Hearts, card.A))
	if got := strategy.Winner(r); got != ErrRoundNotEnded.Error() {
		t.Fatalf("expected %q, got %q", ErrRoundNotEnded.Error(), got)
	}

	// Ended: all hands empty and t2 has higher score.
	r.Teams["t1"].Players[0].Hand = r.Teams["t1"].Players[0].Hand // keep reference explicit
	r.Teams["t1"].Players[0].Hand.Cards = nil
	r.GetScore()["t1"] = 40
	r.GetScore()["t2"] = 80
	if got := strategy.Winner(r); got != "t2" {
		t.Fatalf("expected winner t2, got %q", got)
	}
}

func TestCalculateCurrentTrickRoundPoints(t *testing.T) {
	t.Parallel()

	r := NewRound(uuid.New(), makeRoundTeams(), bot_strategy.NewEasyBotStrategy())
	r.TrumpSuit = card.Spades
	r.CurrentTrick = trick.NewTrick("p1", r.TrumpSuit, r.Teams)

	plays := []trick.Play{
		trick.NewPlay("p1", mkCard(t, card.Clubs, card.A)),
		trick.NewPlay("p2", mkCard(t, card.Spades, card.Two)), // trump winner
		trick.NewPlay("p3", mkCard(t, card.Clubs, card.K)),
		trick.NewPlay("p4", mkCard(t, card.Clubs, card.Q)),
	}
	for _, p := range plays {
		if err := r.CurrentTrick.AddPlay(p); err != nil {
			t.Fatalf("failed to add play %+v: %v", p, err)
		}
	}

	points := r.RuleStrategy.CalculateCurrentTrickRoundPoints(r)
	// A(11) + 2(0) + K(4) + Q(2) = 17 points to t2
	if points["t2"] != 17 {
		t.Fatalf("expected t2 points 17, got %d", points["t2"])
	}
	if points["t1"] != 0 {
		t.Fatalf("expected t1 points 0, got %d", points["t1"])
	}
}
