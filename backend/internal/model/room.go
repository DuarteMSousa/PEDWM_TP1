package domain

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type RoomStatus string

const (
	RoomOpen   RoomStatus = "OPEN"
	RoomInGame RoomStatus = "IN_GAME"
	RoomClosed RoomStatus = "CLOSED"

	SalaAberta  RoomStatus = RoomOpen
	SalaEmJogo  RoomStatus = RoomInGame
	SalaFechada RoomStatus = RoomClosed
)

var (
	ErrInvalidRoomID          = errors.New("invalid room id")
	ErrInvalidHost            = errors.New("invalid host")
	ErrRoomNotOpen            = errors.New("cannot join/leave: room is not open")
	ErrRoomFull               = errors.New("room is full")
	ErrPlayerAlreadyInRoom    = errors.New("player already in room")
	ErrPlayerNotFoundInRoom   = errors.New("player not found in room")
	ErrCannotStartGamePlayers = errors.New("cannot start game: need exactly 4 players")
	ErrInvalidGameID          = errors.New("invalid game id")
)

type Room struct {
	ID      string             `json:"id"`
	HostID  string             `json:"hostId"`
	Players map[string]*Player `json:"players"`
	Status  RoomStatus         `json:"status"`
	GameID  string             `json:"gameId,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewRoom(roomID string, host *Player) (*Room, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}
	if host == nil || strings.TrimSpace(host.ID) == "" {
		return nil, ErrInvalidHost
	}

	players := map[string]*Player{
		host.ID: host,
	}

	return &Room{
		ID:        roomID,
		HostID:    host.ID,
		Players:   players,
		Status:    RoomOpen,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *Room) AddPlayer(player *Player) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != RoomOpen {
		return ErrRoomNotOpen
	}
	if player == nil || strings.TrimSpace(player.ID) == "" {
		return ErrInvalidPlayerID
	}
	if len(r.Players) >= 4 {
		return ErrRoomFull
	}
	if _, exists := r.Players[player.ID]; exists {
		return ErrPlayerAlreadyInRoom
	}

	r.Players[player.ID] = player
	return nil
}

func (r *Room) RemovePlayer(playerID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != RoomOpen {
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
		sort.Strings(ids)
		if len(ids) > 0 {
			r.HostID = ids[0]
		}
	}

	if len(r.Players) == 0 {
		r.Status = RoomClosed
	}

	return nil
}

func (r *Room) CanStartGame() bool {
	if r == nil {
		return false
	}
	return r.Status == RoomOpen && len(r.Players) == 4
}

func (r *Room) StartGame(gameID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if !r.CanStartGame() {
		return ErrCannotStartGamePlayers
	}
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return ErrInvalidGameID
	}

	r.Status = RoomInGame
	r.GameID = gameID
	return nil
}

func (r *Room) Close() {
	if r == nil {
		return
	}
	r.Status = RoomClosed
}
