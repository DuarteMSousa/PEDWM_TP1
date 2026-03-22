package player

import bot_strategy "backend/internal/domain/player/botStrategy"

type Bot struct {
	Player
	Strategy bot_strategy.IBotStrategy
}

func NewBot(id, name string, sequence int, strategy bot_strategy.IBotStrategy) Bot {
	return Bot{
		Player: Player{
			ID:       id,
			Name:     name,
			Type:     BOT,
			Sequence: sequence,
		},
		Strategy: strategy,
	}
}
