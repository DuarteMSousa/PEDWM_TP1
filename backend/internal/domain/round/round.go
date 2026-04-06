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
	ErrInvalidPlay           = errors.New("invalid play")
	ErrTrickNotStarted       = errors.New("trick not started")
)

// Round represents a round of a Sueca game. It manages the deck,
// the trump suit, the tricks, and the accumulated score of each team.
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

// NewRound creates a new round with the specified teams and bot strategy.
func NewRound(gameId uuid.UUID, teams map[string]*team.Team, botStrategy bot_strategy.IBotStrategy) *Round {
	round := &Round{
		gameId:       gameId,
		Teams:        teams,
		BotStrategy:  botStrategy,
		RuleStrategy: NewSuecaRoundRuleStrategy(),
		score:        make(map[string]int),
	}
	for teamID := range teams {
		round.score[teamID] = 0
	}

	round.State = NewRoundSetupState(round)
	return round
}

// StartNewTrick starts a new trick with the specified leader.
func (r *Round) StartNewTrick(leaderID string) {
	r.CurrentTrick = trick.NewTrick(leaderID, r.TrumpSuit, r.Teams)
	r.AddEvent(events.NewTrickStartedEvent(r.gameId.String(), leaderID))
	r.AddEvent(events.NewTurnChangedEvent(r.gameId.String(), leaderID))
	r.State.Update()
}

// GetPlayerTeamId returns the ID of the team to which the player belongs.
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

// GetPlayer returns the player with the specified ID.
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

// PlayCard executes a card play, validating rules and advancing the state.
func (r *Round) PlayCard(playerID string, cardId string) error {
	if r.State == nil {
		return ErrRoundNotStarted
	}
	if r.CurrentTrick == nil {
		return ErrTrickNotStarted
	}

	player, err := r.GetPlayer(playerID)
	if err != nil {
		return err
	}

	card, err := player.Hand.GetCard(cardId)
	if err != nil {
		return err
	}

	play := trick.NewPlay(player.ID, card)
	if !r.CurrentTrick.RuleStrategy.ValidatePlay(*r.CurrentTrick, play) {
		return ErrInvalidPlay
	}

	if _, err := player.Hand.RemoveCard(cardId); err != nil {
		return err
	}

	if err := r.CurrentTrick.AddPlay(play); err != nil {
		player.Hand.AddCard(card)
		return err
	}

	r.AddEvent(events.NewCardPlayedEvent(r.gameId.String(), player.ID, card))

	if r.CurrentTrick != nil && !r.RuleStrategy.HasEnded(r) {
		nextPlayerID, err := r.CurrentTrick.TurnOrder.Next()
		if err == nil {
			r.AddEvent(events.NewTurnChangedEvent(r.gameId.String(), nextPlayerID))
		}
	}

	r.State.Update()

	return nil
}

// GetScore returns the current score of the round by team.
func (r *Round) GetScore() map[string]int {
	return r.score
}

// AddEvent adds an event to the round's event queue.
func (r *Round) AddEvent(event events.Event) events.Event {
	r.events = append(r.events, event)
	return event
}

// CollectEvents returns and clears the accumulated events of the round.
func (r *Round) CollectEvents() []events.Event {
	collectedEvents := r.events
	r.events = []events.Event{}
	return collectedEvents
}
