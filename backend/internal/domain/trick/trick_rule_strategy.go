package trick

type ITrickRuleStrategy interface {
	Winner(trick Trick) (string, error)
	HasEnded(trick Trick) bool
	ValidatePlay(trick Trick, play Play) bool
}
