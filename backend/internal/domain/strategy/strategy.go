package strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/game"
)

// Strategy contracts only.
// Concrete implementations are owned by another teammate.

type TrickRuleStrategy interface {
	Winner(trumpSuit card.Naipe, plays []game.Play) string
}

type ScoringStrategy interface {
	TrickPoints(plays []game.Play) int
}

type BotPlayStrategy interface {
	ChooseCard(hand []card.Card, leadSuit *card.Naipe, trumpSuit card.Naipe) (card.Card, error)
}
