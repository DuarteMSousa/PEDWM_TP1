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

// Team represents a team of players.
type Team struct {
	ID      string
	Players []*player.Player
}

// NewTeam creates a new team, validating that the ID is not empty.
func NewTeam(id string, players []*player.Player) (Team, error) {
	team := Team{
		ID:      id,
		Players: players,
	}

	return team, team.Validate()
}

// Validate checks if the team has a valid ID.
func (t Team) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return ErrInvalidTeamID
	}
	return nil
}
