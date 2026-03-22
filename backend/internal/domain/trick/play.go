package trick

import "backend/internal/domain/card"

type Play struct {
	PlayerID string
	Card     card.Card
}

func NewPlay(playerID string, card card.Card) Play {
	return Play{
		PlayerID: playerID,
		Card:     card,
	}
}
