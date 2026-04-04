package command

import (
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/room"
)

type ChangeBotStrategyCommand struct {
	option bot_strategy.BotStrategyType
}

func NewChangeBotStrategyCommand(option bot_strategy.BotStrategyType) ChangeBotStrategyCommand {
	return ChangeBotStrategyCommand{option: option}
}

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
