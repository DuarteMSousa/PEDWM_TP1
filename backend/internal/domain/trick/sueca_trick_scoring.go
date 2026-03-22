package trick

import "backend/internal/domain/card"

type SuecaTrickScoring struct{}

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

func (s SuecaTrickScoring) TrickPoints(plays []Play) int {

	total := 0
	for _, play := range plays {
		total += s.CardPoints(play.Card)
	}
	return total
}
