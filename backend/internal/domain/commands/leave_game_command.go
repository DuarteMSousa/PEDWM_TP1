package command

import "backend/internal/domain/game"

type LeaveGameCommand struct {
	playerId string
}

func NewLeaveGameCommand(playerId string) LeaveGameCommand {
	return LeaveGameCommand{playerId: playerId}
}

func (c LeaveGameCommand) Execute(game *game.Game) {
	panic("not implemented yet")
	// game.RemovePlayer(c.playerId)
}
