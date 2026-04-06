package round

// IRoundRuleStrategy defines the rules of a round: winner, end, and scoring.
type IRoundRuleStrategy interface {
	Winner(Round *Round) string
	HasEnded(Round *Round) bool
	CalculateCurrentTrickRoundPoints(Round *Round) map[string]int
}
