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
	Hands        map[string]*hand.Hand
	Points       map[string]int
	State        IRoundState
	RuleStrategy IRoundRuleStrategy
}

func NewRound(trump card.Suit, teams []team.Team) *Round {
	if !trump.Valid() {
		panic(card.ErrInvalidSuit)
	}

	players := make([]*player.Player, 0)

	hands := make(map[string]*hand.Hand)
	points := make(map[string]int)
	for _, t := range teams {
		points[t.ID] = 0
		for _, p := range t.Players {
			players = append(players, p)
			hands[p.ID] = hand.NewHand()

		}
	}

	round := &Round{
		TrumpSuit: trump,
		Hands:     hands,
		Points:    points,
	}

	//select random player id
	leaderID := players[rand.Intn(len(players))].ID
	round.StartNewTrick(leaderID)

	round.State = NewRoundSetupState(round)
	return round
}

func (r *Round) StartNewTrick(leaderID string) {
	if r.CurrentTrick == nil {
		r.CurrentTrick = trick.NewTrick(leaderID)
		return
	}
	r.CurrentTrick.Reset(leaderID)
}
