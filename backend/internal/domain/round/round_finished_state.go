package round

import "backend/internal/domain/events"

// RoundFinishedState implements RoundState
type RoundFinishedState struct {
	round *Round
}

// NewRoundFinishedState creates a new instance of RoundFinishedState
func NewRoundFinishedState(round *Round) *RoundFinishedState {
	return &RoundFinishedState{round: round}
}

// Enter is called when the round enters the finished state, it adds a RoundEndedEvent to the event log
func (s *RoundFinishedState) Enter() {
	s.round.AddEvent(events.NewRoundEndedEvent(s.round.gameId.String(), s.round.GetScore(), s.round.RuleStrategy.Winner(s.round)))
}

// Update is called periodically while the round is in the finished state, it does nothing in this state
func (s *RoundFinishedState) Update() {}
