package round

import "errors"

var (
	ErrRoundNotEnded = errors.New("round not ended")
)

// SuecaRoundRuleStrategy implements the rules of Sueca at the round level.
type SuecaRoundRuleStrategy struct{}

// NewSuecaRoundRuleStrategy creates a new instance of the strategy.
func NewSuecaRoundRuleStrategy() *SuecaRoundRuleStrategy {
	return &SuecaRoundRuleStrategy{}
}

// Winner returns the ID of the team that won the round.
func (s *SuecaRoundRuleStrategy) Winner(Round *Round) string {

	if !Round.RuleStrategy.HasEnded(Round) {
		return ErrRoundNotEnded.Error()
	}

	var winner string
	var maxPoints int
	for teamID, _ := range Round.Teams {
		if Round.GetScore()[teamID] > maxPoints {
			maxPoints = Round.GetScore()[teamID]
			winner = teamID
		}
	}

	return winner
}

// HasEnded checks if the round has ended (all hands are empty).
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

// CalculateCurrentTrickRoundPoints calculates the round points
// awarded to each team based on the current trick.
func (s *SuecaRoundRuleStrategy) CalculateCurrentTrickRoundPoints(Round *Round) map[string]int {
	points := make(map[string]int)

	winner, winnerErr := Round.CurrentTrick.RuleStrategy.WinningTeam(*Round.CurrentTrick)

	if winnerErr != nil {
		return points
	}

	for teamID := range Round.Teams {
		teamPoints := 0
		if teamID == winner {
			teamPoints += Round.CurrentTrick.ScoringStrategy.TrickPoints(Round.CurrentTrick.Plays)
		}
		points[teamID] = teamPoints
	}

	return points
}
