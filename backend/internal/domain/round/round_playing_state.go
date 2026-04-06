package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	"errors"
	"log/slog"
	"math/rand"
)

// RoundPlayingState implements RoundState
type RoundPlayingState struct {
	round *Round
}

var (
	ErrHandNotFound      = errors.New("player hand not found")
	ErrBotStrategyNotSet = errors.New("bot strategy not set")
	ErrBotCardNotFound   = errors.New("bot card not found")
)

// NewRoundPlayingState creates a new instance of RoundPlayingState
func NewRoundPlayingState(r *Round) *RoundPlayingState {
	return &RoundPlayingState{round: r}
}

// Enter initializes the round by selecting a random player to start the first trick and then updates the state.
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

// Update checks if the current trick has ended.
// If it has, it calculates the points for the trick, determines the winner, and updates the scores.
// If the round has ended, it transitions to the RoundFinishedState.
// If not, it starts a new trick with the winning player as the leader.
// If the current trick has not ended, it determines the next player and if it's a bot, it uses the bot strategy to choose a card to play.
func (s *RoundPlayingState) Update() {

	if s.round.CurrentTrick.RuleStrategy.HasEnded(*s.round.CurrentTrick) {

		roundPoints := s.round.RuleStrategy.CalculateCurrentTrickRoundPoints(s.round)

		winnerId, err := s.round.CurrentTrick.RuleStrategy.WinningPlayer(*s.round.CurrentTrick)
		winningTeamId, teamErr := s.round.GetPlayerTeamId(winnerId)

		if teamErr != nil {
			panic(teamErr)
		}

		if err != nil {
			panic(err)
		}

		s.round.AddEvent(events.NewTrickEndedEvent(s.round.gameId.String(), winnerId, roundPoints[winningTeamId]))

		for _, team := range s.round.Teams {
			s.round.score[team.ID] += roundPoints[team.ID]
		}

		if s.round.RuleStrategy.HasEnded(s.round) {
			s.round.State = NewRoundFinishedState(s.round)
			s.round.State.Enter()
		} else {
			s.round.StartNewTrick(winnerId)
			s.round.State.Update()
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
			if s.round.BotStrategy == nil {
				panic(ErrBotStrategyNotSet)
			}

			if nextPlayer.Hand == nil {
				panic(ErrHandNotFound)
			}

			leadSuit := card.Suit("")
			if s.round.CurrentTrick.LeadSuit != nil {
				leadSuit = *s.round.CurrentTrick.LeadSuit
			}
			chosenCard := s.round.BotStrategy.ChooseCard(*nextPlayer.Hand, leadSuit, s.round.CurrentTrick.RuleStrategy)
			if chosenCard.ID == "" {
				panic(ErrBotCardNotFound)
			}

			err = s.round.PlayCard(nextPlayer.ID, chosenCard.ID)

			slog.Debug("bot played card", "bot", nextPlayer.Name, "rank", chosenCard.Rank, "suit", chosenCard.Suit)
			if err != nil {
				panic(err)
			}
		}

	}
}
