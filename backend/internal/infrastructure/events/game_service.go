package events_infrastructure

import (
	"backend/internal/domain/game"
)

type GameService interface {
	SetGameStatus(gameID string, status game.GameStatus) (*game.Game, error)
}
