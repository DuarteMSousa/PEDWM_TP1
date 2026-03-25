package room

import (
	"backend/internal/domain/game"
	game_factory "backend/internal/domain/game/gameFactory"
	domainplayer "backend/internal/domain/player"
	"errors"
	"math/rand"
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
)

type Room struct {
	ID      string                          `json:"id"`
	HostID  string                          `json:"hostId"`
	Players map[string]*domainplayer.Player `json:"players"`
	Status  RoomStatus                      `json:"status"`
	Game    *game.Game                      `json:"game,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewRoom(roomID string, host *domainplayer.Player) (*Room, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}
	if host == nil || strings.TrimSpace(host.ID) == "" {
		return nil, ErrInvalidHost
	}

	players := map[string]*domainplayer.Player{
		host.ID: host,
	}

	return &Room{
		ID:        roomID,
		HostID:    host.ID,
		Players:   players,
		Status:    OPEN,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *Room) AddPlayer(p *domainplayer.Player) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != OPEN {
		return ErrRoomNotOpen
	}
	if p == nil || strings.TrimSpace(p.ID) == "" {
		return ErrInvalidPlayerID
	}
	if len(r.Players) >= 4 {
		return ErrRoomFull
	}
	if _, exists := r.Players[p.ID]; exists {
		return ErrPlayerAlreadyInRoom
	}

	r.Players[p.ID] = p
	return nil
}

func (r *Room) RemovePlayer(playerID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != OPEN {
		return ErrRoomNotOpen
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ErrInvalidPlayerID
	}

	if _, exists := r.Players[playerID]; !exists {
		return ErrPlayerNotFoundInRoom
	}

	delete(r.Players, playerID)

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

	if len(r.Players) == 0 {
		r.Status = CLOSED
	}

	return nil
}

func (r *Room) CanStartGame() bool {
	if r == nil {
		return false
	}
	return r.Status == OPEN && len(r.Players) == 4
}

func (r *Room) StartGame() error {

	if !r.CanStartGame() {
		return ErrCannotStartGamePlayers
	}

	r.Game = game_factory.CreateSuecaGame(r.Players)
	r.Status = IN_GAME

	r.Game.State.Enter()

	return nil
}

func (r *Room) Close() {
	if r == nil {
		return
	}
	r.Status = CLOSED
}
