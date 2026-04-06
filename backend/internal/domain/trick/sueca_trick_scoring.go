package trick

import "backend/internal/domain/card"

// SuecaTrickScoring implements the scoring of cards in Sueca.
type SuecaTrickScoring struct{}

// CardPoints returns the points of a card in Sueca (A=11, 7=10, K=4, J=3, Q=2, rest=0).
func (s SuecaTrickScoring) CardPoints(card card.Card) int {
	switch card.Rank {
	case "A":
		return 11
	case "7":
		return 10
	case "K":
		return 4
	case "J":
		return 3
	case "Q":
		return 2
	default:
		return 0
	}
}

// TrickPoints calculates the total points of a trick.
func (s SuecaTrickScoring) TrickPoints(plays []Play) int {

	total := 0
	for _, play := range plays {
		total += s.CardPoints(play.Card)
	}
	return total
}
