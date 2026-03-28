package command

import (
	"backend/internal/domain/game"
	"backend/internal/domain/player"
)

type EnterGameCommand struct {
	player player.Player
}

func NewEnterGameCommand(player player.Player) EnterGameCommand {
	return EnterGameCommand{player: player}
}

func (c EnterGameCommand) Execute(game *game.Game) {
	game.AddPlayer(c.player.ID)
}
