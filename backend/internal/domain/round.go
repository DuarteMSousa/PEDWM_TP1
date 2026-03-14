package domain

import "errors"

var (
	ErrRoundNotStarted = errors.New("round not started")
)

// Round (mão) representa as 10 vazas jogadas com o mesmo baralho (4 jogadores -> 10 vazas).
// Guarda o trunfo, o contador de vazas e a vaza atual.
type Round struct {
	TrumpSuit    Naipe
	TricksPlayed int // 0..10
	CurrentTrick *Trick
}

func NewRound(trump Naipe, firstLeaderID string) (*Round, error) {
	if !trump.Valid() {
		return nil, ErrInvalidNaipe
	}
	return &Round{
		TrumpSuit:    trump,
		TricksPlayed: 0,
		CurrentTrick: NewTrick(firstLeaderID),
	}, nil
}

func (r *Round) IsFinished() bool {
	return r.TricksPlayed >= 10
}

func (r *Round) StartNewTrick(leaderID string) {
	if r.CurrentTrick == nil {
		r.CurrentTrick = NewTrick(leaderID)
		return
	}
	r.CurrentTrick.Reset(leaderID)
}

func (r *Round) IncrementTrick() {
	r.TricksPlayed++
}
