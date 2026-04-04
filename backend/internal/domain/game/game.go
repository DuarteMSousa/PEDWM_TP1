package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/round"
	"backend/internal/domain/team"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	PENDING     GameStatus = "PENDING"
	IN_PROGRESS GameStatus = "IN_PROGRESS"
	FINISHED    GameStatus = "FINISHED"
)

var (
	ErrGameNotPlaying     = errors.New("game not in playing state")
	ErrNotYourTurn        = errors.New("not your turn")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrTeamNotFound       = errors.New("team not found")
	ErrInvalidPlayerOrder = errors.New("invalid player order")
	ErrStrategyNotSet     = errors.New("strategy not set")
	ErrEventBusNotSet     = errors.New("event bus not set")
	ErrRoundNotConfigured = errors.New("round not configured")
	ErrTrickNotConfigured = errors.New("current trick not configured")
	ErrGameDoesNotExist   = errors.New("game does not exist")
)

type Game struct {
	ID      uuid.UUID
	RoomID  string
	Status  GameStatus
	Players map[string]*player.Player
	Teams   map[string]*team.Team
	Score   map[string]int

	State           IGameState
	round           *round.Round
	scoringStrategy IGameScoringStrategy
	botStrategy     bot_strategy.IBotStrategy

	Events               []events.Event
	currentEventSequence int

	eventBus *events.EventBus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGame(teams []*team.Team, scoringStrategy IGameScoringStrategy, botStrategy bot_strategy.IBotStrategy) *Game {
	g := &Game{
		ID:                   uuid.New(),
		Status:               PENDING,
		Players:              make(map[string]*player.Player),
		Score:                make(map[string]int),
		scoringStrategy:      scoringStrategy,
		botStrategy:          botStrategy,
		Teams:                make(map[string]*team.Team),
		Events:               []events.Event{},
		currentEventSequence: 1,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	for _, t := range teams {
		for _, p := range t.Players {
			g.Players[p.ID] = p
		}
	}
	g.scoringStrategy = scoringStrategy
	g.botStrategy = botStrategy
	g.Teams = make(map[string]*team.Team)
	for _, t := range teams {
		g.Teams[t.ID] = t
	}
	g.Events = []events.Event{}
	g.eventBus = events.DefaultBus()

	g.State = NewGameStartingState(g)
	return g
}

func (g *Game) AddEvent(e events.Event) {
	if g == nil {
		return
	}

	if g.ID != uuid.Nil {
		e.GameID = g.ID.String()
	}

	roomID := strings.TrimSpace(g.RoomID)
	if roomID != "" {
		e.RoomID = roomID
	}

	e.Sequence = g.currentEventSequence
	g.currentEventSequence++

	g.Events = append(g.Events, e)
	if g.eventBus != nil {
		g.eventBus.Publish(e)
	}

}

func (g *Game) GetEvents() []events.Event {
	return g.Events
}

func (g *Game) PlayCard(playerId string, cardId string) error {
	if g == nil {
		return ErrGameDoesNotExist
	}
	if g.State == nil {
		return ErrGameNotPlaying
	}
	if g.round == nil {
		return ErrRoundNotConfigured
	}

	_, ok := g.Players[playerId]
	if !ok {
		return ErrPlayerNotFound
	}

	if err := g.round.PlayCard(playerId, cardId); err != nil {
		return err
	}

	g.State.Update()
	return nil
}

func (g *Game) UpdateRoundState() {
	g.round.State.Update()
	events := g.round.CollectEvents()
	for _, event := range events {
		event.GameID = g.ID.String()
		event.RoomID = g.RoomID

		g.AddEvent(event)
	}
}

func (g *Game) SetEventBus(eventBus *events.EventBus) {
	g.eventBus = eventBus
}

func (g *Game) CurrentRound() *round.Round {
	if g == nil {
		return nil
	}
	return g.round
}

func (g *Game) GetPlayer(playerID string) (*player.Player, error) {
	if g == nil {
		return nil, ErrGameDoesNotExist
	}
	p, ok := g.Players[playerID]
	if !ok {
		return nil, ErrPlayerNotFound
	}
	return p, nil
}
