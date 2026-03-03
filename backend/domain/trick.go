package domain

import "errors"

var (
	ErrTrickComplete     = errors.New("trick already complete")
	ErrPlayerAlreadyPlay = errors.New("player already played in this trick")
)

// Trick representa uma vaza em curso (até 4 jogadas).
// Ela guarda as jogadas e o naipe de saída (lead suit).
type Trick struct {
	LeaderID string
	LeadSuit *Naipe
	Plays    []Play
}

func NewTrick(leaderID string) *Trick {
	return &Trick{
		LeaderID: leaderID,
		Plays:    make([]Play, 0, 4),
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

// AddPlay adiciona a jogada à vaza. Não valida regras de Sueca (seguir naipe/trunfo);
// isso deve ser feito pela TrickRuleStrategy (ValidatePlay / Winner).
func (t *Trick) AddPlay(play Play) error {
	if t.IsComplete() {
		return ErrTrickComplete
	}
	if t.HasPlayed(play.PlayerID) {
		return ErrPlayerAlreadyPlay
	}

	if t.IsEmpty() {
		ls := play.Card.Naipe
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