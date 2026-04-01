package trick

import "backend/internal/domain/card"

type ITrickRuleStrategy interface {
	WinningTeam(trick Trick) (string, error)
	WinningPlayer(trick Trick) (string, error)
	HasEnded(trick Trick) bool
	ValidatePlay(trick Trick, play Play) bool
	CardStrength(card card.Rank) int
}
