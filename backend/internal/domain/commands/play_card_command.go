package command

import "backend/internal/domain/game"

// PlayCardCommand encapsulates the action of playing a card.
type PlayCardCommand struct {
	playerId string
	cardId   string
}

// NewPlayCardCommand creates a command to play a card.
func NewPlayCardCommand(playerId string, cardId string) PlayCardCommand {
	return PlayCardCommand{playerId: playerId, cardId: cardId}
}

// Execute executes the command to play a card on the game.
func (c PlayCardCommand) Execute(game *game.Game) error {
	return game.PlayCard(c.playerId, c.cardId)
}
