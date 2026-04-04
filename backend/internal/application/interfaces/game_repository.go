package interfaces

import "backend/internal/domain/game"

type GameRepository interface {
	Save(g *game.Game) error
	FindByID(id string) (*game.Game, error)
	FindByRoomID(roomID string) ([]*game.Game, error)
	GetByUserID(userID string) ([]*game.Game, error)
}
