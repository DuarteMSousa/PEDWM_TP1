package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/deck"
	"backend/internal/domain/hand"
	"backend/internal/domain/player"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"errors"
	"math/rand"
)

var (
	ErrRoundNotStarted = errors.New("round not started")
)

type Round struct {
	TrumpSuit    card.Suit
	CurrentTrick *trick.Trick
	Deck         *deck.Deck
	Teams        map[string]team.Team
	State        IRoundState
	RuleStrategy IRoundRuleStrategy
}

func NewRound(trump card.Suit, teams []team.Team) *Round {
	if !trump.Valid() {
		panic(card.ErrInvalidSuit)
	}

	players := make([]*player.Player, 0)

	teamsMap := make(map[string]team.Team)
	for _, t := range teams {
		t.RoundScore = 0
		for _, p := range t.Players {
			players = append(players, p)
			p.Hand = hand.NewHand()
		}
		teamsMap[t.ID] = t
	}

	round := &Round{
		TrumpSuit: trump,
		Teams:     teamsMap,
	}

	for _, t := range teams {
		round.Teams[t.ID] = t

	}

	//select random player id
	leaderID := players[rand.Intn(len(players))].ID
	round.StartNewTrick(leaderID)

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
