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
	players map[string]*player.Player
	Teams   map[string]*team.Team
	Score   map[string]int

	State           IGameState
	round           *round.Round
	scoringStrategy IGameScoringStrategy
	botStrategy     bot_strategy.IBotStrategy

	events   []events.Event
	eventBus *events.EventBus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGame(teams []*team.Team, scoringStrategy IGameScoringStrategy, botStrategy bot_strategy.IBotStrategy) *Game {
	g := &Game{
		ID:              uuid.New(),
		Status:          IN_PROGRESS,
		players:         make(map[string]*player.Player),
		Score:           make(map[string]int),
		scoringStrategy: scoringStrategy,
		botStrategy:     botStrategy,
		Teams:           make(map[string]*team.Team),
		events:          []events.Event{},
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	for _, t := range teams {
		for _, p := range t.Players {
			g.players[p.ID] = p
		}
	}
	g.scoringStrategy = scoringStrategy
	g.botStrategy = botStrategy
	g.Teams = make(map[string]*team.Team)
	for _, t := range teams {
		g.Teams[t.ID] = t
	}
	g.events = []events.Event{}
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

	g.events = append(g.events, e)
	if g.eventBus != nil {
		g.eventBus.Publish(e)
	}

}

func (g *Game) GetEvents() []events.Event {
	return g.events
}

func (g *Game) PlayCard(playerId string, cardId string) {
	if g.State == nil {
		panic(ErrGameNotPlaying)
	}

	_, ok := g.players[playerId]
	if !ok {
		panic(ErrPlayerNotFound)
	}

	g.round.PlayCard(playerId, cardId)

	g.State.Update()

}

func (g *Game) UpdateRoundState() {
	g.round.State.Update()
	events := g.round.CollectEvents()
	for _, event := range events {
		g.AddEvent(event)
	}
}

// func (g *Game) SwitchPlayer(player player.Player) error {
// 	if g == nil {
// 		return ErrGameDoesNotExist
// 	}

// 	for _, qplayer := range g.players {
// 		if player.Sequence == qplayer.Sequence {
// 			g.players[qplayer.ID] = &player
// 		}
// 	}

// 	return nil
// }
