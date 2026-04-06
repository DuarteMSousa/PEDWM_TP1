package command

import (
	"backend/internal/domain/game"
)

// ICommand defines an interface for executable commands on a game.
type ICommand interface {
	Execute(game *game.Game) error
}
