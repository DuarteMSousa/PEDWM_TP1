package domain

import (
	"errors"
	"strings"
)

type PlayerType string

const (
	Humano PlayerType = "HUMANO"
	Bot    PlayerType = "BOT"
)

var (
	ErrCardNotFound    = errors.New("card not found")
	ErrInvalidPlayer   = errors.New("invalid player")
	ErrInvalidPlayerID = errors.New("invalid player id")
)

type Player struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Type   PlayerType `json:"type"`
	TeamID string     `json:"teamId,omitempty"`
	Hand   []Card     `json:"hand,omitempty"`
}

func (p Player) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return ErrInvalidPlayerID
	}
	if strings.TrimSpace(p.Name) == "" {
		return ErrInvalidPlayer
	}
	return nil
}

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
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return c, true
		}
	}
	return Card{}, false
}
