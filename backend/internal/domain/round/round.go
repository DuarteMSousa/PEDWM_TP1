package round

import (
	"backend/internal/domain/card"
	"backend/internal/domain/deck"
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrRoundNotStarted       = errors.New("round not started")
	ErrRoundFinished         = errors.New("round finished")
	ErrWinningPlayerNotFound = errors.New("winning player not found in any team")
	ErrPlayerNotFound        = errors.New("player not found in any team")
)

type Round struct {
	gameId       uuid.UUID
	TrumpSuit    card.Suit
	CurrentTrick *trick.Trick
	Deck         *deck.Deck
	Teams        map[string]*team.Team
	State        IRoundState
	RuleStrategy IRoundRuleStrategy
	BotStrategy  bot_strategy.IBotStrategy
	score        map[string]int

	events []events.Event
}

func NewRound(gameId uuid.UUID, teams map[string]*team.Team, botStrategy bot_strategy.IBotStrategy) *Round {
	round := &Round{
		gameId:      gameId,
		Teams:       teams,
		BotStrategy: botStrategy,
	}

	round.State = NewRoundSetupState(round)
	return round
}

func (r *Round) StartNewTrick(leaderID string) {
	r.CurrentTrick = trick.NewTrick(leaderID, r.TrumpSuit, r.Teams)
	r.AddEvent(events.NewTrickStartedEvent(r.gameId.String(), leaderID))
}

func (r *Round) GetPlayerTeamId(playerID string) (string, error) {
	for _, team := range r.Teams {
		for _, player := range team.Players {
			if player.ID == playerID {
				return team.ID, nil
			}
		}
	}
	return "", ErrWinningPlayerNotFound
}

func (r *Round) GetPlayer(playerID string) (*player.Player, error) {
	for _, team := range r.Teams {
		for _, player := range team.Players {
			if player.ID == playerID {
				return player, nil
			}
		}
	}
	return nil, ErrPlayerNotFound
}

func (r *Round) PlayCard(playerID string, cardId string) error {
	if r.State == nil {
		panic(ErrRoundNotStarted)
	}

	player, err := r.GetPlayer(playerID)
	if err != nil {
		return err
	}

	card, err := player.Hand.RemoveCard(cardId)
	if err != nil {
		return err
	}

	play := trick.NewPlay(player.ID, card)

	r.CurrentTrick.AddPlay(play)
	r.AddEvent(events.NewCardPlayedEvent(r.CurrentTrick.LeaderID, player.ID, card))

	r.State.Update()

	return nil
}

func (r *Round) GetScore() map[string]int {
	return r.score
}

func (r *Round) AddEvent(event events.Event) events.Event {
	r.events = append(r.events, event)
	return event
}

func (r *Round) CollectEvents() []events.Event {
	collectedEvents := r.events
	r.events = []events.Event{}
	return collectedEvents
}
