package round

import (
	"backend/internal/domain/player"
	"math/rand"
)

// RoundPlayingState implementa RoundState
type RoundPlayingState struct {
	round *Round
}

func NewRoundPlayingState(r *Round) *RoundPlayingState {
	return &RoundPlayingState{round: r}
}

func (s *RoundPlayingState) Enter() {
	players := make([]*player.Player, 0)

	for _, team := range s.round.Teams {
		for _, player := range team.Players {
			players = append(players, player)
		}
	}

	firstLeaderId := players[rand.Intn(len(players))].ID

	s.round.StartNewTrick(firstLeaderId)

	s.round.State.Update()
}
func (s *RoundPlayingState) Update() {
	if s.round.RuleStrategy.HasEnded(s.round) {
		s.round.State = NewRoundFinishedState(s.round)
		s.round.State.Enter()
		return
	}

	if s.round.CurrentTrick.RuleStrategy.HasEnded(*s.round.CurrentTrick) {
		//Fazer distribuição de pontos e continuar a ronda
	}

}
