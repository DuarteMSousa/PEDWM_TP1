package command

import (
	"backend/internal/domain/game"
)

type ICommand interface {
	Execute(game *game.Game) error
}
