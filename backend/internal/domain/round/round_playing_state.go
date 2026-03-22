package round

// RoundPlayingState implementa RoundState
type RoundPlayingState struct {
	round *Round
}

func NewRoundPlayingState(r *Round) *RoundPlayingState {
	return &RoundPlayingState{round: r}
}

func (s *RoundPlayingState) Enter() {

}
func (s *RoundPlayingState) Update() {

}
