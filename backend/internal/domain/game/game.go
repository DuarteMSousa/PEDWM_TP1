package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	"backend/internal/domain/round"
	"backend/internal/domain/strategy"
	"backend/internal/domain/team"
	"errors"
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
	ID      string
	players map[string]*player.Player
	teams   map[string]*team.Team

	state           IGameState
	round           *round.Round
	scoringStrategy strategy.ScoringStrategy
	botStrategy     strategy.BotPlayStrategy

	events   []*events.Event
	eventBus *events.EventBus
}

func (g *Game) NewGame(teams []team.Team, scoringStrategy strategy.ScoringStrategy, botStrategy strategy.BotPlayStrategy) {
	g.ID = "game-123" // Gerar ID único
	g.players = make(map[string]*player.Player)
	for _, t := range teams {
		for _, p := range t.Players {
			g.players[p.ID] = &p
		}
	}
	g.scoringStrategy = scoringStrategy
	g.botStrategy = botStrategy
	g.teams = make(map[string]*team.Team)
	for _, t := range teams {
		g.teams[t.ID] = &t
	}
	g.events = []*events.Event{}
	g.state = NewGameStartingState(g)
}
