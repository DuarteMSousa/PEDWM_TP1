package round

import "errors"

var (
	ErrRoundNotEnded = errors.New("round not ended")
)

type SuecaRoundRuleStrategy struct{}

func NewSuecaRoundRuleStrategy() *SuecaRoundRuleStrategy {
	return &SuecaRoundRuleStrategy{}
}

func (s *SuecaRoundRuleStrategy) Winner(Round *Round) string {

	if !Round.RuleStrategy.HasEnded(Round) {
		return ErrRoundNotEnded.Error()
	}

	var winner string
	var maxPoints int
	for teamID, team := range Round.Teams {
		if team.RoundScore > maxPoints {
			maxPoints = team.RoundScore
			winner = teamID
		}
	}

	return winner
}

func (s *SuecaRoundRuleStrategy) HasEnded(Round *Round) bool {
	endedHands := 0
	totalPlayers := 0
	for _, team := range Round.Teams {
		for _, player := range team.Players {
			if player.Hand.IsEmpty() {
				endedHands++
			}
			totalPlayers++
		}
	}

	return endedHands == totalPlayers
}
