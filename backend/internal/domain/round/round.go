package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/deck"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"errors"
)

var (
	ErrRoundNotStarted = errors.New("round not started")
)

type Round struct {
	TrumpSuit    card.Suit
	CurrentTrick *trick.Trick
	Deck         *deck.Deck
	Teams        map[string]*team.Team
	State        IRoundState
	RuleStrategy IRoundRuleStrategy
}

func NewRound(teams map[string]*team.Team) *Round {
	round := &Round{
		Teams: teams,
	}

	round.State = NewRoundSetupState(round)
	return round
}

func (r *Round) StartNewTrick(leaderID string) {
	if r.CurrentTrick == nil {
		r.CurrentTrick = trick.NewTrick(leaderID, r.TrumpSuit, r.Teams)
		return
	}
	r.CurrentTrick.Reset(leaderID)
}
