package trick

import (
	"backend/internal/domain/card"
	"errors"
)

// SuecaTrickRules implements the rules of Sueca for a trick.
type SuecaTrickRules struct{}

var (
	ErrTrickNotEnded         = errors.New("trick has not ended yet")
	ErrWinningPlayerNotFound = errors.New("winning player not found in any team")
)

// WinningPlayer determines the player who won the trick.
// The winner is determined by the highest card of the leading suit,
// or by a trump card if present.
func (s SuecaTrickRules) WinningPlayer(trick Trick) (string, error) {
	if !trick.RuleStrategy.HasEnded(trick) {
		return "", ErrTrickNotEnded
	}

	winningPlay := trick.Plays[0]

	for _, play := range trick.Plays[1:] {
		if play.Card.Suit == winningPlay.Card.Suit {
			if s.CardStrength(play.Card.Rank) > s.CardStrength(winningPlay.Card.Rank) {
				winningPlay = play
			}
		} else if play.Card.Suit == trick.TrumpSuit && winningPlay.Card.Suit != trick.TrumpSuit {
			winningPlay = play
		}
	}

	return winningPlay.PlayerID, nil
}

// CardStrength returns the strength of a card in Sueca (A=10, 7=9, K=8, ...).
func (s SuecaTrickRules) CardStrength(r card.Rank) int {
	switch r {
	case card.A:
		return 10
	case card.Seven:
		return 9
	case card.K:
		return 8
	case card.J:
		return 7
	case card.Q:
		return 6
	case card.Six:
		return 5
	case card.Five:
		return 4
	case card.Four:
		return 3
	case card.Three:
		return 2
	case card.Two:
		return 1
	default:
		return 0
	}
}

// WinningTeam returns the ID of the team that won the trick.
func (s SuecaTrickRules) WinningTeam(trick Trick) (string, error) {
	winningPlayerID, err := s.WinningPlayer(trick)
	if err != nil {
		return "", err
	}

	for teamID, team := range trick.Teams {
		for _, player := range team.Players {
			if player.ID == winningPlayerID {
				return teamID, nil
			}
		}
	}

	return "", ErrWinningPlayerNotFound

}

// HasEnded checks if all players have played in this trick.
func (s SuecaTrickRules) HasEnded(trick Trick) bool {
	playerCount := 0
	for _, team := range trick.Teams {
		playerCount += len(team.Players)
	}

	return len(trick.Plays) == playerCount
}

// ValidatePlay checks if the play respects the rules of Sueca
// (turn order, obligation to follow the leading suit).
func (s SuecaTrickRules) ValidatePlay(trick Trick, play Play) bool {

	if trick.RuleStrategy.HasEnded(trick) {
		return false
	}

	nextPlayerID, err := trick.TurnOrder.Next()

	if err != nil {
		return false
	}

	if nextPlayerID != play.PlayerID {
		return false
	}

	if trick.LeadSuit == nil {
		return true
	}

	if play.Card.Suit != *trick.LeadSuit {

		for _, t := range trick.Teams {
			for _, p := range t.Players {
				if p.ID == play.PlayerID {
					if p.Hand.HasSuit(*trick.LeadSuit) {
						return false
					}
				}
			}
		}
	}

	return true

}
