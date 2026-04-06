package player

import "backend/internal/domain/hand"

// PlayerType distinguishes between human players and bots.
type PlayerType string

const (
	HUMAN PlayerType = "HUMAN"
	BOT   PlayerType = "BOT"
)

// Player represents a player in a room or game.
type Player struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     PlayerType `json:"type"`
	Sequence int        `json:"sequence"`
	Hand     *hand.Hand `json:"hand"`
}

// NewPlayer creates a human player with an empty hand.
func NewPlayer(id, name string, sequence int) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Type:     HUMAN,
		Sequence: sequence,
		Hand:     hand.NewHand(),
	}
}
