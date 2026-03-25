package round

type IRoundRuleStrategy interface {
	Winner(Round *Round) string
	HasEnded(Round *Round) bool
	CalculateCurrentTrickRoundPoints(Round *Round) map[string]int
}
