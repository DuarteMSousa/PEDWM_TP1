package command

import (
	"backend/internal/domain/game"
	"errors"
)

type LeaveGameCommand struct {
	playerId string
}

func NewLeaveGameCommand(playerId string) LeaveGameCommand {
	return LeaveGameCommand{playerId: playerId}
}

func (c LeaveGameCommand) Execute(game *game.Game) error {
	return errors.New("leave game command not implemented yet")
	// game.RemovePlayer(c.playerId)
}
