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
	ErrPlayerOutOfTurn   = errors.New("player is playing out of turn")
)

// Trick represents a trick in progress (up to 4 plays).
// It keeps track of the plays and the leading suit.
type Trick struct {
	LeaderID  string
	LeadSuit  *card.Suit
	TrumpSuit card.Suit
	Plays     []Play
	Teams     map[string]*team.Team

	TurnOrder       turnorder.TurnOrder
	ScoringStrategy ITrickScoringStrategy
	RuleStrategy    ITrickRuleStrategy
}

// NewTrick creates a new trick with the given leader, trump suit, and teams.
// It initializes the turn order starting from the leader.
func NewTrick(leaderID string, TrumpSuit card.Suit, teams map[string]*team.Team) *Trick {
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
		LeaderID:        leaderID,
		TurnOrder:       turnorder,
		TrumpSuit:       TrumpSuit,
		Plays:           make([]Play, 0),
		Teams:           teams,
		ScoringStrategy: SuecaTrickScoring{},
		RuleStrategy:    SuecaTrickRules{},
	}
}

// IsEmpty indicates if no plays have been made in this trick.
func (t *Trick) IsEmpty() bool {
	return len(t.Plays) == 0
}

// IsComplete indicates if all players have played.
func (t *Trick) IsComplete() bool {
	numPlayers := 0
	for _, team := range t.Teams {
		numPlayers += len(team.Players)
	}

	return len(t.Plays) == numPlayers
}

// HasPlayed indicates if the player has already played in this trick.
func (t *Trick) HasPlayed(playerID string) bool {
	for _, p := range t.Plays {
		if p.PlayerID == playerID {
			return true
		}
	}
	return false
}

// AddPlay adds a play to the trick after validating turn order and completeness.
func (t *Trick) AddPlay(play Play) error {
	if t.IsComplete() {
		return ErrTrickComplete
	}

	if t.HasPlayed(play.PlayerID) {
		return ErrPlayerAlreadyPlay
	}

	nextPlayerID, err := t.TurnOrder.Next()

	if err != nil {
		return err
	}

	if nextPlayerID != play.PlayerID {
		return ErrPlayerOutOfTurn
	}

	if t.IsEmpty() {
		ls := play.Card.Suit
		t.LeadSuit = &ls
	}

	t.Plays = append(t.Plays, play)
	_, err = t.TurnOrder.Dequeue()
	if err != nil {
		return err
	}

	return nil
}

// Reset resets the trick for a new leader, clearing the plays.
func (t *Trick) Reset(newLeaderID string) {
	t.LeaderID = newLeaderID
	t.LeadSuit = nil
	t.Plays = t.Plays[:0]
}
