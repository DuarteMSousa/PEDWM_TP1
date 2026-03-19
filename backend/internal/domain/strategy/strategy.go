package strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/trick"
)

// Strategy contracts only.
// Concrete implementations are owned by another teammate.

type TrickRuleStrategy interface {
	Winner(trumpSuit card.Naipe, plays []trick.Play) string
}

type ScoringStrategy interface {
	TrickPoints(plays []trick.Play) int
}

type BotPlayStrategy interface {
	ChooseCard(hand []card.Card, leadSuit *card.Naipe, trumpSuit card.Naipe) (card.Card, error)
}
