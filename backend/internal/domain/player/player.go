package player

type PlayerType string

const (
	HUMAN PlayerType = "HUMAN"
	BOT   PlayerType = "BOT"
)

type Player struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Type PlayerType `json:"type"`
}

type Bot struct {
	Player,
	Strategy BotStrategy
}

func NewPlayer(id, name string, playerType PlayerType) Player {
	return Player{
		ID:   id,
		Name: name,
		Type: playerType,
	}
}

func NewBot(id, name string, strategy BotStrategy) Bot {
	return Bot{
		Player:   NewPlayer(id, name, BOT),
		Strategy: strategy,
	}
}
