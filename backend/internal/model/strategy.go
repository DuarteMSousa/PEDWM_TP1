package domain

// Strategy contracts only.
// Concrete implementations are owned by another teammate.

type TrickRuleStrategy interface {
	Winner(trumpSuit Naipe, plays []Play) string
}

type ScoringStrategy interface {
	TrickPoints(plays []Play) int
}

type BotPlayStrategy interface {
	ChooseCard(hand []Card, leadSuit *Naipe, trumpSuit Naipe) (Card, error)
}
