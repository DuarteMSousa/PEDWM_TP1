package command

import "backend/internal/domain/game"

type PlayCardCommand struct {
	playerId string
	cardId   string
}

func NewPlayCardCommand(playerId string, cardId string) PlayCardCommand {
	return PlayCardCommand{playerId: playerId, cardId: cardId}
}

func (c PlayCardCommand) Execute(game *game.Game) error {
	return game.PlayCard(c.playerId, c.cardId)
}
