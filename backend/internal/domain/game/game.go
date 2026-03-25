package game

import (
	"backend/internal/domain/events"
	game_strategy "backend/internal/domain/game/gameStrategy"
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
	teams   map[string]*team.Team

	state           IGameState
	round           *round.Round
	scoringStrategy game_strategy.IGameScoringStrategy
	botStrategy     bot_strategy.IBotStrategy

	events   []*events.Event
	eventBus *events.EventBus
}

func NewGame(teams []*team.Team, scoringStrategy game_strategy.IGameScoringStrategy, botStrategy bot_strategy.IBotStrategy) *Game {
	g := &Game{
		ID:              uuid.New(),
		players:         make(map[string]*player.Player),
		scoringStrategy: scoringStrategy,
		botStrategy:     botStrategy,
		teams:           make(map[string]*team.Team),
		events:          []*events.Event{},
	}

	for _, t := range teams {
		for _, p := range t.Players {
			g.players[p.ID] = p
		}
	}
	g.scoringStrategy = scoringStrategy
	g.botStrategy = botStrategy
	g.teams = make(map[string]*team.Team)
	for _, t := range teams {
		g.teams[t.ID] = t
	}
	g.events = []*events.Event{}
	g.eventBus = events.NewEventBus()

	g.state = NewGameStartingState(g)
	g.state.Enter()
	return g
}
