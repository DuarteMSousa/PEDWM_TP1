package game

import (
	"backend/internal/domain/events"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/round"
	"backend/internal/domain/team"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type gameObserverStub struct {
	events []events.Event
}

func (o *gameObserverStub) Update(event events.Event) {
	o.events = append(o.events, event)
}

type fakeGameState struct{}

func (fakeGameState) Enter()  {}
func (fakeGameState) Update() {}

func makeBaseTeams() []*team.Team {
	p1 := player.NewPlayer("p1", "P1", 1)
	p2 := player.NewPlayer("p2", "P2", 2)
	p3 := player.NewPlayer("p3", "P3", 3)
	p4 := player.NewPlayer("p4", "P4", 4)

	t1 := &team.Team{ID: "t1", Players: []*player.Player{p1, p3}}
	t2 := &team.Team{ID: "t2", Players: []*player.Player{p2, p4}}
	return []*team.Team{t1, t2}
}

func makeGame() *Game {
	return NewGame(
		makeBaseTeams(),
		NewSuecaGameScoringStrategy(),
		bot_strategy.NewEasyBotStrategy(),
	)
}

func TestNewGameInitialStateAndMaps(t *testing.T) {
	t.Parallel()

	g := makeGame()
	if g.Status != PENDING {
		t.Fatalf("expected status PENDING, got %s", g.Status)
	}
	if len(g.Players) != 4 {
		t.Fatalf("expected 4 players in map, got %d", len(g.Players))
	}
	if len(g.Teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(g.Teams))
	}
	if g.State == nil {
		t.Fatal("expected initial game state to be set")
	}
}

func TestAddEventFillsGameAndRoomAndSequenceAndPublishes(t *testing.T) {
	t.Parallel()

	g := makeGame()
	g.RoomID = "room_1"

	bus := events.NewEventBus()
	obs := &gameObserverStub{}
	bus.Subscribe(obs)
	g.SetEventBus(bus)

	g.AddEvent(events.NewTurnChangedEvent("", "p1"))
	g.AddEvent(events.NewRoundStartedEvent(""))

	if len(g.Events) != 2 {
		t.Fatalf("expected 2 game events, got %d", len(g.Events))
	}

	first := g.Events[0]
	if first.GameID != g.ID.String() {
		t.Fatalf("expected GameID %q, got %q", g.ID.String(), first.GameID)
	}
	if first.RoomID != "room_1" {
		t.Fatalf("expected RoomID room_1, got %q", first.RoomID)
	}
	if first.Sequence != 1 {
		t.Fatalf("expected first sequence 1, got %d", first.Sequence)
	}
	if g.Events[1].Sequence != 2 {
		t.Fatalf("expected second sequence 2, got %d", g.Events[1].Sequence)
	}

	if len(obs.events) != 2 {
		t.Fatalf("expected observer to receive 2 events, got %d", len(obs.events))
	}
}

func TestPlayCardErrorGuards(t *testing.T) {
	t.Parallel()

	var nilGame *Game
	err := nilGame.PlayCard("p1", "A_HEARTS")
	if !errors.Is(err, ErrGameDoesNotExist) {
		t.Fatalf("expected ErrGameDoesNotExist, got %v", err)
	}

	g := &Game{}
	err = g.PlayCard("p1", "A_HEARTS")
	if !errors.Is(err, ErrGameNotPlaying) {
		t.Fatalf("expected ErrGameNotPlaying, got %v", err)
	}
}

func TestPlayCardReturnsPlayerNotFoundBeforeRoundExecution(t *testing.T) {
	t.Parallel()

	g := makeGame()
	g.State = fakeGameState{}
	g.round = &round.Round{}

	err := g.PlayCard("missing", "A_HEARTS")
	if !errors.Is(err, ErrPlayerNotFound) {
		t.Fatalf("expected ErrPlayerNotFound, got %v", err)
	}
}

func TestRemovePlayerReplacesWithBotAndAddsEvents(t *testing.T) {
	t.Parallel()

	g := makeGame()
	g.State = fakeGameState{}

	err := g.RemovePlayer("p1")
	if err != nil {
		t.Fatalf("expected RemovePlayer to succeed, got %v", err)
	}

	if _, ok := g.Players["p1"]; ok {
		t.Fatal("expected removed player p1 not to exist")
	}

	botID := "bP1"
	bot, ok := g.Players[botID]
	if !ok {
		t.Fatalf("expected replacement bot %q to exist", botID)
	}
	if bot.Type != player.BOT {
		t.Fatalf("expected replacement player type BOT, got %s", bot.Type)
	}

	if len(g.Events) < 2 {
		t.Fatalf("expected at least 2 events (PLAYER_LEFT + PLAYER_JOINED), got %d", len(g.Events))
	}
	if g.Events[len(g.Events)-2].Type != events.EventPlayerLeft {
		t.Fatalf("expected penultimate event PLAYER_LEFT, got %s", g.Events[len(g.Events)-2].Type)
	}
	if g.Events[len(g.Events)-1].Type != events.EventPlayerJoined {
		t.Fatalf("expected last event PLAYER_JOINED, got %s", g.Events[len(g.Events)-1].Type)
	}
}

func TestGetPlayerAndGetPlayerTeamErrors(t *testing.T) {
	t.Parallel()

	g := makeGame()

	if _, err := g.GetPlayer("missing"); !errors.Is(err, ErrPlayerNotFound) {
		t.Fatalf("expected ErrPlayerNotFound, got %v", err)
	}
	if _, err := g.GetPlayerTeam("missing"); !errors.Is(err, ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got %v", err)
	}
}

func TestSuecaGameScoringStrategyThresholdsAndWinner(t *testing.T) {
	t.Parallel()

	teams := map[string]*team.Team{
		"t1": {ID: "t1", Players: []*player.Player{
			player.NewPlayer("p1", "P1", 1),
			player.NewPlayer("p3", "P3", 3),
		}},
		"t2": {ID: "t2", Players: []*player.Player{
			player.NewPlayer("p2", "P2", 2),
			player.NewPlayer("p4", "P4", 4),
		}},
	}
	r := round.NewRound(uuid.New(), teams, bot_strategy.NewEasyBotStrategy())
	r.GetScore()["t1"] = 95
	r.GetScore()["t2"] = 59

	strategy := NewSuecaGameScoringStrategy()
	points := strategy.CalculateCurrentRoundGamePoints(r)
	if points["t1"] != 2 {
		t.Fatalf("expected t1 to receive 2 game points, got %d", points["t1"])
	}
	if points["t2"] != 0 {
		t.Fatalf("expected t2 to receive 0 game points, got %d", points["t2"])
	}

	g := makeGame()
	g.Score["t1"] = 4
	g.Score["t2"] = 2
	if !strategy.HasGameEnded(g) {
		t.Fatal("expected game to be ended when one team reaches 4 points")
	}
	if winner := strategy.Winner(g); winner != "t1" {
		t.Fatalf("expected winner t1, got %q", winner)
	}
}
