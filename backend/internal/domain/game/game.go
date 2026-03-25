package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/round"
	"backend/internal/domain/team"
	"errors"

	"github.com/google/uuid"
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
)

type Game struct {
	ID      uuid.UUID
	players map[string]*player.Player
	Teams   map[string]*team.Team

	State           IGameState
	round           *round.Round
	scoringStrategy IGameScoringStrategy
	botStrategy     bot_strategy.IBotStrategy

	events   []events.Event
	eventBus *events.EventBus
}

func NewGame(teams []*team.Team, scoringStrategy IGameScoringStrategy, botStrategy bot_strategy.IBotStrategy) *Game {
	g := &Game{
		ID:              uuid.New(),
		players:         make(map[string]*player.Player),
		scoringStrategy: scoringStrategy,
		botStrategy:     botStrategy,
		Teams:           make(map[string]*team.Team),
		events:          []events.Event{},
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
	g.eventBus = events.NewEventBus()

	g.State = NewGameStartingState(g)
	return g
}

func (g *Game) AddEvent(e events.Event) {
	if g == nil {
		return
	}

	g.events = append(g.events, e)
	if g.eventBus != nil {
		g.eventBus.Publish(e)
	}

}
