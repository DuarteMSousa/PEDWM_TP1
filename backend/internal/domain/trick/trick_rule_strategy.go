package trick

import "backend/internal/domain/card"

// ITrickRuleStrategy defines the rules of a trick: winner, validation, and end.
type ITrickRuleStrategy interface {
	WinningTeam(trick Trick) (string, error)
	WinningPlayer(trick Trick) (string, error)
	HasEnded(trick Trick) bool
	ValidatePlay(trick Trick, play Play) bool
	CardStrength(card card.Rank) int
}
