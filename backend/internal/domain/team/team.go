package team

import (
	"backend/internal/domain/player"
	"errors"
	"strings"
)

var (
	ErrInvalidTeamID  = errors.New("invalid team id")
	ErrNegativePoints = errors.New("points cannot be negative")
)

type Team struct {
	ID      string
	Players []*player.Player
}

func NewTeam(id string, players []*player.Player) (Team, error) {
	team := Team{
		ID:      id,
		Players: players,
	}

	return team, team.Validate()
}

func (t Team) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return ErrInvalidTeamID
	}
	return nil
}
