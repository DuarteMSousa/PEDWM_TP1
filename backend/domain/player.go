package domain

import (
	"errors"
	"strings"
)

type PlayerType string

const (
	Humano PlayerType = "HUMANO"
	Bot   PlayerType = "BOT"
)

var (
	ErrCardNotFound   = errors.New("card not found")
	ErrInvalidPlayer  = errors.New("invalid player")
	ErrInvalidPlayerID = errors.New("invalid player id")
)

// Player representa um jogador no domínio.
type Player struct {
	ID     string
	Name   string
	Type   PlayerType
	TeamID string
	Hand   []Card
}

// Validate valida invariantes básicas do jogador.
// (Ajusta conforme o teu domínio: por exemplo, Name pode ser opcional.)
func (p Player) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return ErrInvalidPlayerID
	}
	if strings.TrimSpace(p.TeamID) == "" {
		return ErrInvalidPlayer
	}
	return nil
}

// RemoveCard remove uma carta da mão pelo ID.
// Retorna a carta removida e true se existir; caso contrário Card{} e false.
func (p *Player) RemoveCard(cardID string) (Card, bool) {
	if p == nil {
		return Card{}, false
	}
	cardID = strings.TrimSpace(cardID)
	if cardID == "" {
		return Card{}, false
	}

	for i, c := range p.Hand {
		if c.ID == cardID {
			// Remove mantendo a ordem relativa das restantes
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return c, true
		}
	}
	return Card{}, false
}