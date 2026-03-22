package player

import "backend/internal/domain/hand"

type PlayerType string

const (
	HUMAN PlayerType = "HUMAN"
	BOT   PlayerType = "BOT"
)

type Player struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     PlayerType `json:"type"`
	Sequence int        `json:"sequence"`
	Hand     *hand.Hand `json:"hand"`
}

func NewPlayer(id, name string, sequence int) Player {
	return Player{
		ID:       id,
		Name:     name,
		Type:     HUMAN,
		Sequence: sequence,
		Hand:     hand.NewHand(),
	}
}
