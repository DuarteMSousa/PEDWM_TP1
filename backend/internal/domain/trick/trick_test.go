package trick

import (
	"backend/internal/domain/card"
	"backend/internal/domain/player"
	"backend/internal/domain/team"
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

func makeTeams() map[string]*team.Team {
	p1 := player.NewPlayer("p1", "P1", 1)
	p2 := player.NewPlayer("p2", "P2", 2)
	p3 := player.NewPlayer("p3", "P3", 3)
	p4 := player.NewPlayer("p4", "P4", 4)

	return map[string]*team.Team{
		"t1": {ID: "t1", Players: []*player.Player{p1, p3}},
		"t2": {ID: "t2", Players: []*player.Player{p2, p4}},
	}
}

func TestTrickAddPlaySetsLeadSuitAndEnforcesTurnOrder(t *testing.T) {
	t.Parallel()

	teams := makeTeams()
	tr := NewTrick("p1", card.Spades, teams)

	first := NewPlay("p1", makeCard(t, card.Clubs, card.A))
	if err := tr.AddPlay(first); err != nil {
		t.Fatalf("expected first play to be valid, got %v", err)
	}
	if tr.LeadSuit == nil || *tr.LeadSuit != card.Clubs {
		t.Fatalf("expected lead suit CLUBS, got %v", tr.LeadSuit)
	}

	outOfTurn := NewPlay("p3", makeCard(t, card.Clubs, card.K))
	err := tr.AddPlay(outOfTurn)
	if !errors.Is(err, ErrPlayerOutOfTurn) {
		t.Fatalf("expected ErrPlayerOutOfTurn, got %v", err)
	}
}

func TestTrickAddPlayRejectsDuplicatePlayer(t *testing.T) {
	t.Parallel()

	teams := makeTeams()
	tr := NewTrick("p1", card.Spades, teams)

	first := NewPlay("p1", makeCard(t, card.Hearts, card.A))
	if err := tr.AddPlay(first); err != nil {
		t.Fatalf("expected first play to be valid, got %v", err)
	}

	duplicate := NewPlay("p1", makeCard(t, card.Hearts, card.K))
	err := tr.AddPlay(duplicate)
	if !errors.Is(err, ErrPlayerAlreadyPlay) {
		t.Fatalf("expected ErrPlayerAlreadyPlay, got %v", err)
	}
}

func TestValidatePlayRequiresFollowingLeadSuit(t *testing.T) {
	t.Parallel()

	teams := makeTeams()
	teams["t2"].Players[0].Hand.AddCard(makeCard(t, card.Hearts, card.Two))
	teams["t2"].Players[0].Hand.AddCard(makeCard(t, card.Clubs, card.K))

	tr := NewTrick("p1", card.Spades, teams)
	if err := tr.AddPlay(NewPlay("p1", makeCard(t, card.Hearts, card.A))); err != nil {
		t.Fatalf("failed to add opening play: %v", err)
	}

	invalid := NewPlay("p2", makeCard(t, card.Clubs, card.K))
	if tr.RuleStrategy.ValidatePlay(*tr, invalid) {
		t.Fatal("expected play to be invalid because player must follow lead suit")
	}
}

func TestValidatePlayAllowsDifferentSuitWhenPlayerDoesNotHaveLeadSuit(t *testing.T) {
	t.Parallel()

	teams := makeTeams()
	teams["t2"].Players[0].Hand.AddCard(makeCard(t, card.Clubs, card.K))

	tr := NewTrick("p1", card.Spades, teams)
	if err := tr.AddPlay(NewPlay("p1", makeCard(t, card.Hearts, card.A))); err != nil {
		t.Fatalf("failed to add opening play: %v", err)
	}

	valid := NewPlay("p2", makeCard(t, card.Clubs, card.K))
	if !tr.RuleStrategy.ValidatePlay(*tr, valid) {
		t.Fatal("expected play to be valid when player has no lead suit")
	}
}

func TestWinningPlayerAndTeamTrumpBeatsLeadSuit(t *testing.T) {
	t.Parallel()

	teams := makeTeams()
	tr := NewTrick("p1", card.Spades, teams)

	plays := []Play{
		NewPlay("p1", makeCard(t, card.Clubs, card.A)),
		NewPlay("p2", makeCard(t, card.Spades, card.Two)), // trump
		NewPlay("p3", makeCard(t, card.Clubs, card.K)),
		NewPlay("p4", makeCard(t, card.Clubs, card.Q)),
	}
	for _, p := range plays {
		if err := tr.AddPlay(p); err != nil {
			t.Fatalf("unexpected error adding play %+v: %v", p, err)
		}
	}

	winnerPlayer, err := tr.RuleStrategy.WinningPlayer(*tr)
	if err != nil {
		t.Fatalf("expected winner player, got error %v", err)
	}
	if winnerPlayer != "p2" {
		t.Fatalf("expected winner p2, got %q", winnerPlayer)
	}

	winnerTeam, err := tr.RuleStrategy.WinningTeam(*tr)
	if err != nil {
		t.Fatalf("expected winner team, got error %v", err)
	}
	if winnerTeam != "t2" {
		t.Fatalf("expected winner team t2, got %q", winnerTeam)
	}
}

func TestSuecaTrickScoringPoints(t *testing.T) {
	t.Parallel()

	scoring := SuecaTrickScoring{}
	plays := []Play{
		NewPlay("p1", makeCard(t, card.Hearts, card.A)),     // 11
		NewPlay("p2", makeCard(t, card.Spades, card.Seven)), // 10
		NewPlay("p3", makeCard(t, card.Clubs, card.K)),      // 4
		NewPlay("p4", makeCard(t, card.Diamonds, card.Two)), // 0
	}

	got := scoring.TrickPoints(plays)
	if got != 25 {
		t.Fatalf("expected trick points 25, got %d", got)
	}
}
