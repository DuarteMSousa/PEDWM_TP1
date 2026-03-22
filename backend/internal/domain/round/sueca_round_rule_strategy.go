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
	for teamID, points := range Round.Points {
		if points > maxPoints {
			maxPoints = points
			winner = teamID
		}
	}

	return winner
}

func (s *SuecaRoundRuleStrategy) HasEnded(Round *Round) bool {
	endedHands := 0
	for _, hand := range Round.Hands {
		if hand.IsEmpty() {
			endedHands++
		}
	}

	return endedHands == len(Round.Hands)
}
