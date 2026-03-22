package strategy

import (
	"backend/internal/domain/card"
	"backend/internal/domain/trick"
)

// Strategy contracts only.
// Concrete implementations are owned by another teammate.

type TrickRuleStrategy interface {
	Winner(trumpSuit card.Suit, plays []trick.Play) string
}

// type ScoringStrategy interface {
// 	TrickPoints(plays []trick.Play) int
// }

// type BotPlayStrategy interface {
// 	ChooseCard(hand []card.Card, leadSuit *card.Suit, trumpSuit card.Suit) (card.Card, error)
// }
