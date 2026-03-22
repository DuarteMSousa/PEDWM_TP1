package trick

import (
	"backend/internal/domain/card"
	"backend/internal/domain/player"
	"backend/internal/domain/team"
	"backend/internal/domain/turnorder"
	"errors"
)

var (
	ErrTrickComplete     = errors.New("trick already complete")
	ErrPlayerAlreadyPlay = errors.New("player already played in this trick")
)

// Trick representa uma vaza em curso (até 4 jogadas).
// Ela guarda as jogadas e o Suit de saída (lead suit).
type Trick struct {
	LeaderID  string
	LeadSuit  *card.Suit
	TrumpSuit card.Suit
	Plays     []Play
	Teams     map[string]team.Team

	TurnOrder       turnorder.TurnOrder
	ScoringStrategy ITrickScoringStrategy
	RuleStrategy    ITrickRuleStrategy
}

func NewTrick(leaderID string, TrumpSuit card.Suit, teams map[string]team.Team) *Trick {
	players := make([]*player.Player, 0)
	for _, t := range teams {
		for _, p := range t.Players {
			players = append(players, p)
		}
	}
	turnorder, err := turnorder.NewTurnOrder(leaderID, players)

	if err != nil {
		panic("Failed to create turn order: " + err.Error())
	}

	return &Trick{
		LeaderID:  leaderID,
		TurnOrder: turnorder,
		TrumpSuit: TrumpSuit,
		Plays:     make([]Play, 0),
		Teams:     teams,
	}
}

func (t *Trick) IsEmpty() bool {
	return len(t.Plays) == 0
}

func (t *Trick) IsComplete() bool {
	return len(t.Plays) >= 4
}

func (t *Trick) HasPlayed(playerID string) bool {
	for _, p := range t.Plays {
		if p.PlayerID == playerID {
			return true
		}
	}
	return false
}

func (t *Trick) AddPlay(play Play) error {
	if t.IsComplete() {
		return ErrTrickComplete
	}
	if t.HasPlayed(play.PlayerID) {
		return ErrPlayerAlreadyPlay
	}

	if t.IsEmpty() {
		ls := play.Card.Suit
		t.LeadSuit = &ls
	}

	t.Plays = append(t.Plays, play)
	return nil
}

func (t *Trick) Reset(newLeaderID string) {
	t.LeaderID = newLeaderID
	t.LeadSuit = nil
	t.Plays = t.Plays[:0]
}
