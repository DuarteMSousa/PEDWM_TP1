package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/deck"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"errors"
)

var (
	ErrRoundNotStarted       = errors.New("round not started")
	ErrRoundFinished         = errors.New("round finished")
	ErrWinningPlayerNotFound = errors.New("winning player not found in any team")
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
	r.CurrentTrick = trick.NewTrick(leaderID, r.TrumpSuit, r.Teams)
}

func (r *Round) GetPlayerTeamId(playerID string) (string, error) {
	for _, team := range r.Teams {
		for _, player := range team.Players {
			if player.ID == playerID {
				return team.ID, nil
			}
		}
	}
	return "", ErrWinningPlayerNotFound
}
