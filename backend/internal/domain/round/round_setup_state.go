package round

// RoundSetupState implementa GameState
type RoundSetupState struct {
	round *Round
}

func NewRoundSetupState(r *Round) *RoundSetupState {
	return &RoundSetupState{round: r}
}

func (s *RoundSetupState) Enter() {

}

func (s *RoundSetupState) Update() {

}
