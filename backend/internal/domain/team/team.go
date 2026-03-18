package team

import (
	"errors"
	"strings"
)

var (
	ErrInvalidTeamID  = errors.New("invalid team id")
	ErrNegativePoints = errors.New("points cannot be negative")
)

// Team representa uma equipa no domínio.
type Team struct {
	ID      string
	Players []string // IDs dos jogadores
	Score   int
}

// Validate valida invariantes básicas da equipa.
func (t Team) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return ErrInvalidTeamID
	}
	return nil
}

// AddPoints adiciona pontos à equipa.
// Nota: se no teu jogo pode haver pontos negativos, remove esta validação.
func (t *Team) AddPoints(points int) error {
	if t == nil {
		return ErrInvalidTeamID
	}
	if points < 0 {
		return ErrNegativePoints
	}
	t.Score += points
	return nil
}
