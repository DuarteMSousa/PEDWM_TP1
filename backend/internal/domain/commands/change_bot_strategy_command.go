package command

import (
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/room"
)

// ChangeBotStrategyCommand encapsulates the action of changing the bots' strategy.
type ChangeBotStrategyCommand struct {
	option bot_strategy.BotStrategyType
}

// NewChangeBotStrategyCommand creates a command to change the bot strategy.
func NewChangeBotStrategyCommand(option bot_strategy.BotStrategyType) ChangeBotStrategyCommand {
	return ChangeBotStrategyCommand{option: option}
}

// Execute executes the command to change the bot strategy on the room.
func (c ChangeBotStrategyCommand) Execute(room *room.Room) error {
	var botstrategy bot_strategy.IBotStrategy = bot_strategy.NewEasyBotStrategy()

	switch c.option {
	case bot_strategy.EASY:
		botstrategy = bot_strategy.NewEasyBotStrategy()
		break
	case bot_strategy.HARD:
		botstrategy = bot_strategy.NewHardBotStrategy()
		break
	}

	room.SetBotStrategy(botstrategy)
	return nil
}
