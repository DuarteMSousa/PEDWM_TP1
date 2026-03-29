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
	for teamID, _ := range Round.Teams {
		if Round.GetScore()[teamID] > maxPoints {
			maxPoints = Round.GetScore()[teamID]
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

func (s *SuecaRoundRuleStrategy) CalculateCurrentTrickRoundPoints(Round *Round) (map[string]int, error) {
	points := make(map[string]int)

	winner, winnerErr := Round.CurrentTrick.RuleStrategy.WinningTeam(*Round.CurrentTrick)

	if winnerErr != nil {
		return nil, winnerErr
	}

	for teamID := range Round.Teams {
		teamPoints := 0
		if teamID == winner {
			teamPoints += Round.CurrentTrick.ScoringStrategy.TrickPoints(Round.CurrentTrick.Plays)
		}
		points[teamID] = teamPoints
	}

	return points, nil
}
