package command

import (
	"backend/internal/domain/room"
)

// ICommand defines an interface for executable commands on a game.
type ICommand interface {
	Execute(game *room.Room) error
}
