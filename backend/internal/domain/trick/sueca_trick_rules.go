package trick

import (
	"backend/internal/domain/card"
	"errors"
)

type SuecaTrickRules struct{}

var (
	ErrTrickNotEnded         = errors.New("trick has not ended yet")
	ErrWinningPlayerNotFound = errors.New("winning player not found in any team")
)

func (s SuecaTrickRules) WinningPlayer(trick Trick) (string, error) {
	if !trick.RuleStrategy.HasEnded(trick) {
		return "", ErrTrickNotEnded
	}

	winningPlay := trick.Plays[0]

	for _, play := range trick.Plays[1:] {
		if play.Card.Suit == winningPlay.Card.Suit {
			if s.TrickStrength(play.Card.Rank) > s.TrickStrength(winningPlay.Card.Rank) {
				winningPlay = play
			}
		} else if play.Card.Suit == trick.TrumpSuit && winningPlay.Card.Suit != trick.TrumpSuit {
			winningPlay = play
		}
	}

	return winningPlay.PlayerID, nil
}

func (s SuecaTrickRules) TrickStrength(r card.Rank) int {
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

func (s SuecaTrickRules) HasEnded(trick Trick) bool {
	playerCount := 0
	for _, team := range trick.Teams {
		playerCount += len(team.Players)
	}

	return len(trick.Plays) == playerCount
}

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

		if play.Card.Suit != trick.TrumpSuit {
			for _, t := range trick.Teams {
				for _, p := range t.Players {
					if p.ID == play.PlayerID {
						if p.Hand.HasSuit(*trick.LeadSuit) || p.Hand.HasSuit(trick.TrumpSuit) {
							return false
						}
					}
				}
			}
		} else {
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

	}

	return true

}
