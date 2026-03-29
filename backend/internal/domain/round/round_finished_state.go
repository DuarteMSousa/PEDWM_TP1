package round

import "backend/internal/domain/events"

// RoundFinishedState implementa RoundState
type RoundFinishedState struct {
	round *Round
}

func NewRoundFinishedState(round *Round) *RoundFinishedState {
	return &RoundFinishedState{round: round}
}

func (s *RoundFinishedState) Enter() {
	s.round.AddEvent(events.NewRoundEndedEvent(s.round.gameId.String(), s.round.GetScore(), s.round.RuleStrategy.Winner(s.round)))
}

func (s *RoundFinishedState) Update() {}
