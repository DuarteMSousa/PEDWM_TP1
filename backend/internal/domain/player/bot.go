package player

import bot_strategy "backend/internal/domain/player/botStrategy"

// Bot is a Player with an automatic game strategy.
type Bot struct {
	Player
	Strategy bot_strategy.IBotStrategy
}

// NewBot creates a bot with the specified strategy.
func NewBot(id, name string, sequence int, strategy bot_strategy.IBotStrategy) *Bot {
	return &Bot{
		Player: Player{
			ID:       id,
			Name:     name,
			Type:     BOT,
			Sequence: sequence,
		},
		Strategy: strategy,
	}
}
