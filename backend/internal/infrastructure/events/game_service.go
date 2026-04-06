package events_infrastructure

import (
	"backend/internal/domain/game"
)

// GameService defines the interface for interacting with the game service.
type GameService interface {
	SetGameStatus(gameID string, status game.GameStatus) (*game.Game, error)
}
