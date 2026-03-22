package player

type PlayerType string

const (
	HUMAN PlayerType = "HUMAN"
	BOT   PlayerType = "BOT"
)

type Player struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     PlayerType `json:"type"`
	Sequence int        `json:"sequence"`
}

func NewPlayer(id, name string, sequence int) Player {
	return Player{
		ID:       id,
		Name:     name,
		Type:     HUMAN,
		Sequence: sequence,
	}
}
