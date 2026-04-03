package room

import (
	"backend/internal/domain/events"
	"backend/internal/domain/game"
	game_factory "backend/internal/domain/game/gameFactory"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"errors"
	"sort"
	"strings"
	"time"
)

type RoomStatus string

const (
	OPEN    RoomStatus = "OPEN"
	IN_GAME RoomStatus = "IN_GAME"
	CLOSED  RoomStatus = "CLOSED"
)

var (
	ErrInvalidRoomID          = errors.New("invalid room id")
	ErrInvalidHost            = errors.New("invalid room host")
	ErrRoomNotOpen            = errors.New("room is not open")
	ErrRoomFull               = errors.New("room is full")
	ErrPlayerAlreadyInRoom    = errors.New("player already in room")
	ErrPlayerNotFoundInRoom   = errors.New("player not found in room")
	ErrCannotStartGamePlayers = errors.New("cannot start game: need exactly 4 players")
	ErrInvalidGameID          = errors.New("invalid game id")
	ErrInvalidPlayerID        = errors.New("invalid player id")
	ErrGameInProgress         = errors.New("cannot leave room while game is in progress")
)

type Room struct {
	ID       string                    `json:"id"`
	HostID   string                    `json:"hostId"`
	Players  map[string]*player.Player `json:"players"`
	Status   RoomStatus                `json:"status"`
	Game     *game.Game                `json:"game,omitempty"`
	EventBus *events.EventBus          `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewRoom(id string, hostId string, hostUsername string) (*Room, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalidRoomID
	}

	hostId = strings.TrimSpace(hostId)
	if hostId == "" {
		return nil, ErrInvalidHost
	}
	if hostUsername == "" {
		return nil, ErrInvalidHost
	}

	players := map[string]*player.Player{
		hostId: player.NewPlayer(hostId, hostUsername, 1),
	}

	return &Room{
		ID:        id,
		HostID:    hostId,
		Players:   players,
		Status:    OPEN,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *Room) AddPlayer(userID, username string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != OPEN {
		return ErrRoomNotOpen
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidPlayerID
	}
	if len(r.Players) >= 4 {
		return ErrRoomFull
	}
	if _, exists := r.Players[userID]; exists {
		return ErrPlayerAlreadyInRoom
	}

	joinedPlayer := player.NewPlayer(userID, username, len(r.Players)+1)
	r.Players[userID] = joinedPlayer

	if r.EventBus != nil {
		event := events.NewPlayerJoinedEvent(
			"",
			joinedPlayer.ID,
			joinedPlayer.Name,
			joinedPlayer.Sequence,
		)
		event.RoomID = r.ID
		r.EventBus.Publish(event)
	}

	return nil
}

func (r *Room) RemovePlayer(playerID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	r.SyncStatusWithGame()

	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ErrInvalidPlayerID
	}
	if r.Status == IN_GAME {
		return ErrGameInProgress
	}
	if r.Status != OPEN {
		return ErrRoomNotOpen
	}

	if _, exists := r.Players[playerID]; !exists {
		return ErrPlayerNotFoundInRoom
	}

	removedPlayer := r.Players[playerID]
	delete(r.Players, playerID)

	if r.EventBus != nil {
		gameID := ""
		if r.Game != nil {
			gameID = r.Game.ID.String()
		}
		event := events.NewPlayerLeftEvent(gameID, removedPlayer.ID, r.ID)
		event.RoomID = r.ID
		r.EventBus.Publish(event)
	}

	if r.HostID == playerID {
		nextHostID := r.findNextHostID()
		if nextHostID != "" {
			r.HostID = nextHostID
		}
	}

	r.reindexSequences()

	if len(r.Players) == 0 {
		r.Status = CLOSED
		if r.EventBus != nil {
			event := events.NewRoomClosedEvent(r.ID)
			event.RoomID = r.ID
			r.EventBus.Publish(event)
		}
	}

	return nil
}

func (r *Room) CanStartGame() bool {
	if r == nil {
		return false
	}
	r.SyncStatusWithGame()
	return r.Status == OPEN
}

func (r *Room) StartGame() error {
	r.SyncStatusWithGame()
	if !r.CanStartGame() {
		return ErrCannotStartGamePlayers
	}

	gamePlayers := make(map[string]*player.Player, len(r.Players))
	for id, p := range r.Players {
		gamePlayers[id] = p
	}

	r.Game = game_factory.CreateSuecaGame(gamePlayers, bot_strategy.NewEasyBotStrategy(), r.EventBus)
	r.Game.RoomID = r.ID
	r.Status = IN_GAME

	r.Game.State.Enter()

	return nil
}

func (r *Room) SetEventBus(eventBus *events.EventBus) {
	r.EventBus = eventBus
}

func (r *Room) Close() {
	if r == nil {
		return
	}
	r.Status = CLOSED
}

func (r *Room) SyncStatusWithGame() bool {
	if r == nil || r.Game == nil {
		return false
	}
	if r.Status == IN_GAME && r.Game.Status == game.FINISHED {
		r.Status = OPEN
		return true
	}
	return false
}

func (r *Room) reindexSequences() {
	if r == nil || len(r.Players) == 0 {
		return
	}

	orderedPlayers := make([]*player.Player, 0, len(r.Players))
	for _, p := range r.Players {
		orderedPlayers = append(orderedPlayers, p)
	}
	sort.Slice(orderedPlayers, func(i, j int) bool {
		if orderedPlayers[i].Sequence == orderedPlayers[j].Sequence {
			return orderedPlayers[i].ID < orderedPlayers[j].ID
		}
		return orderedPlayers[i].Sequence < orderedPlayers[j].Sequence
	})

	for index, p := range orderedPlayers {
		p.Sequence = index + 1
	}
}

func (r *Room) findNextHostID() string {
	if r == nil || len(r.Players) == 0 {
		return ""
	}

	orderedPlayers := make([]*player.Player, 0, len(r.Players))
	for _, p := range r.Players {
		orderedPlayers = append(orderedPlayers, p)
	}
	sort.Slice(orderedPlayers, func(i, j int) bool {
		if orderedPlayers[i].Sequence == orderedPlayers[j].Sequence {
			return orderedPlayers[i].ID < orderedPlayers[j].ID
		}
		return orderedPlayers[i].Sequence < orderedPlayers[j].Sequence
	})

	return orderedPlayers[0].ID
}
