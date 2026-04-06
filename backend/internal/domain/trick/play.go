package trick

import "backend/internal/domain/card"

// Play represents a player's move (played card).
type Play struct {
	PlayerID string
	Card     card.Card
}

// NewPlay creates a new play.
func NewPlay(playerID string, card card.Card) Play {
	return Play{
		PlayerID: playerID,
		Card:     card,
	}
}
