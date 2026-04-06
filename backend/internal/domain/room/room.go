package room

import (
	"backend/internal/domain/events"
	"backend/internal/domain/game"
	game_factory "backend/internal/domain/game/gameFactory"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"errors"
	"log/slog"
	"math/rand"
	"strings"
	"time"
)

// RoomStatus represents the status of a room.
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
)

// Room is the root aggregate that represents a game room.
type Room struct {
	ID       string                    `json:"id"`
	HostID   string                    `json:"hostId"`
	Players  map[string]*player.Player `json:"players"`
	Status   RoomStatus                `json:"status"`
	Game     *game.Game                `json:"game,omitempty"`
	EventBus *events.EventBus          `json:"-"`

	BotStrategy bot_strategy.IBotStrategy `json:"-"`
	CreatedAt   time.Time                 `json:"createdAt"`
}

// NewRoom creates a new open room with the specified host.
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
		ID:          id,
		HostID:      hostId,
		Players:     players,
		Status:      OPEN,
		CreatedAt:   time.Now().UTC(),
		BotStrategy: bot_strategy.NewEasyBotStrategy(),
	}, nil
}

// AddPlayer adds a player to the room (maximum 4). Publishes PlayerJoinedEvent.
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

// RemovePlayer removes a player from the room. If the game is active,
// the player is replaced by a bot. If the room becomes empty, it is closed.
func (r *Room) RemovePlayer(playerID string) error {
	slog.Debug("removing a player", "playerID", playerID, "roomID", r.ID)
	if r == nil {
		return errors.New("room is nil")
	}

	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ErrInvalidPlayerID
	}

	if _, exists := r.Players[playerID]; !exists {
		return ErrPlayerNotFoundInRoom
	}

	removedPlayer := r.Players[playerID]
	delete(r.Players, playerID)

	if r.Game != nil {
		playerErr := r.Game.RemovePlayer(playerID)

		if playerErr != nil {
			slog.Warn("failed to synchronize game state after player left", "playerID", playerID, "roomID", r.ID, "error", playerErr)
		}
	} else {
		if r.EventBus != nil {
			gameID := ""
			event := events.NewPlayerLeftEvent(gameID, removedPlayer.ID, r.ID)
			event.RoomID = r.ID
			r.EventBus.Publish(event)
		}
	}

	if r.HostID == playerID {
		r.HostID = ""
		ids := make([]string, 0, len(r.Players))
		for id := range r.Players {
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			r.HostID = ids[rand.Intn(len(ids))]
		}
	}

	slog.Info("player left the room", "playerID", removedPlayer.ID, "roomID", r.ID)
	slog.Debug("remaining players in the room", "roomID", r.ID, "count", len(r.Players))

	if len(r.Players) == 0 {
		r.Status = CLOSED
		if r.EventBus != nil {
			event := events.NewRoomClosedEvent(r.ID)
			event.RoomID = r.ID
			r.EventBus.Publish(event)
			slog.Info("ROOM_CLOSED event published", "roomID", r.ID)
		}
	}

	return nil
}

// CanStartGame checks if the room is in the OPEN state.
func (r *Room) CanStartGame() bool {
	if r == nil {
		return false
	}
	return r.Status == OPEN
}

// CreateGame creates a Sueca game from the players in the room.
func (r *Room) CreateGame() error {
	if !r.CanStartGame() {
		return ErrCannotStartGamePlayers
	}

	gamePlayers := make(map[string]*player.Player, len(r.Players))
	for id, p := range r.Players {
		gamePlayers[id] = p
	}

	r.Game = game_factory.CreateSuecaGame(gamePlayers, r.BotStrategy, r.EventBus)
	r.Game.RoomID = r.ID
	r.Status = IN_GAME

	return nil
}

// SetEventBus sets the event bus for publishing room events.
func (r *Room) SetEventBus(eventBus *events.EventBus) {
	r.EventBus = eventBus
}

// Close closes the room by setting its status to CLOSED.
func (r *Room) Close() {
	if r == nil {
		return
	}
	r.Status = CLOSED
}

// SetBotStrategy sets the bot strategy and publishes BotStrategyChangedEvent.
func (r *Room) SetBotStrategy(strategy bot_strategy.IBotStrategy) {
	if r == nil {
		return
	}
	r.BotStrategy = strategy
	if r.EventBus != nil {
		event := events.NewBotStrategyChangedEvent(r.ID, strategy.GetType())
		r.EventBus.Publish(event)
	}
}
