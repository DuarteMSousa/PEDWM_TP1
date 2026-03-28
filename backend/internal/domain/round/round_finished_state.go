package round

// RoundFinishedState implementa RoundState
type RoundFinishedState struct {
	round *Round
}

func NewRoundFinishedState(round *Round) *RoundFinishedState {
	return &RoundFinishedState{round: round}
}

func (s *RoundFinishedState) Enter() {}

func (s *RoundFinishedState) Update() {}
