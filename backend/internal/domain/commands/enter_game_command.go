package command

import (
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	"errors"
)

type EnterGameCommand struct {
	player player.Player
}

func NewEnterGameCommand(player player.Player) EnterGameCommand {
	return EnterGameCommand{player: player}
}

func (c EnterGameCommand) Execute(game *game.Game) error {
	return errors.New("enter game command not implemented yet")
	// game.AddPlayer(c.player.ID)
}
