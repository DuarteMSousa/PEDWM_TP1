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

	if s.round.CurrentTrick.RuleStrategy.HasEnded(*s.round.CurrentTrick) {
		roundPoints := s.round.RuleStrategy.CalculateCurrentTrickRoundPoints(s.round)

		winnerId, err := s.round.CurrentTrick.RuleStrategy.WinningPlayer(*s.round.CurrentTrick)

		if err != nil {
			panic(err)
		}

		for _, team := range s.round.Teams {
			team.RoundScore += roundPoints[team.ID]
		}

		if s.round.RuleStrategy.HasEnded(s.round) {
			s.round.State = NewRoundFinishedState(s.round)
			s.round.State.Enter()
		} else {
			s.round.StartNewTrick(winnerId)
		}
	} else {
		nextId, err := s.round.CurrentTrick.TurnOrder.Next()

		if err != nil {
			panic(err)
		}

		nextPlayer, playerErr := s.round.GetPlayer(nextId)

		if playerErr != nil {
			panic(playerErr)
		}

		if nextPlayer.Type == player.BOT {
			// choosenCard := s.round.BotStrategy.ChooseCard(*nextPlayer.Hand, s.round.TrumpSuit)

			// s.round.PlayCard(*nextPlayer, choosenCard.ID)
		}

	}
}
