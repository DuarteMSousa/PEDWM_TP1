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
	ID         string
	Players    []*player.Player
	GameScore  int
	RoundScore int
}

func NewTeam(id string, players []*player.Player) (Team, error) {
	team := Team{
		ID:         id,
		Players:    players,
		GameScore:  0,
		RoundScore: 0,
	}

	return team, team.Validate()
}

func (t Team) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return ErrInvalidTeamID
	}
	return nil
}

func (t *Team) AddRoundPoints(points int) error {
	if t == nil {
		return ErrInvalidTeamID
	}
	if points < 0 {
		return ErrNegativePoints
	}
	t.RoundScore += points
	return nil
}

func (t *Team) AddGamePoints(points int) error {
	if t == nil {
		return ErrInvalidTeamID
	}
	if points < 0 {
		return ErrNegativePoints
	}
	t.GameScore += points
	return nil
}
