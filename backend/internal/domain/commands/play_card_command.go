package command

import (
	"backend/internal/domain/room"
	"errors"
)

var (
	ErrNoActiveGame = errors.New("No game found in this room")
)

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
func (c PlayCardCommand) Execute(room *room.Room) error {

	if room.Game == nil {
		return ErrNoActiveGame
	}

	return room.Game.PlayCard(c.playerId, c.cardId)
}
