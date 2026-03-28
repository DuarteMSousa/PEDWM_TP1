package command

import "backend/internal/domain/game"

type LeaveGameCommand struct {
	playerId string
}

func NewLeaveGameCommand(playerId string) LeaveGameCommand {
	return LeaveGameCommand{playerId: playerId}
}

func (c LeaveGameCommand) Execute(game *game.Game) {
	game.RemovePlayer(c.playerId)
}
