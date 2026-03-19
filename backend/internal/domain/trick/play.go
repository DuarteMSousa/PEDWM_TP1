package trick

import "backend/internal/domain/card"

type Play struct {
	PlayerID string
	Card     card.Card
}
