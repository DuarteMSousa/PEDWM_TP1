package lobby

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DefaultMaxPlayers = 4
)

type RoomStatus string

const (
	RoomStatusOpen   RoomStatus = "OPEN"
	RoomStatusInGame RoomStatus = "IN_GAME"
	RoomStatusClosed RoomStatus = "CLOSED"
)

var (
	ErrNicknameRequired = errors.New("nickname is required")
	ErrInvalidPlayerID  = errors.New("player id is required")
	ErrPlayerNotFound   = errors.New("player not found")

	ErrInvalidRoomID    = errors.New("room id is required")
	ErrRoomNameRequired = errors.New("room name is required")
	ErrRoomNotFound     = errors.New("room not found")
	ErrRoomNotOpen      = errors.New("cannot join/leave: room is not open")
	ErrRoomFull         = errors.New("room is full")
	ErrNotRoomHost      = errors.New("only the room host can delete this room")
)

type Player struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"createdAt"`
}

type Room struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	HostID     string     `json:"hostId"`
	Status     RoomStatus `json:"status"`
	MaxPlayers int        `json:"maxPlayers"`
	IsPrivate  bool       `json:"isPrivate"`
	PlayerIDs  []string   `json:"playerIds"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type RoomView struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PlayersCount int    `json:"playersCount"`
	MaxPlayers   int    `json:"maxPlayers"`
	IsPrivate    bool   `json:"isPrivate"`
}

type roomRecord struct {
	ID         string
	Name       string
	HostID     string
	Status     RoomStatus
	MaxPlayers int
	IsPrivate  bool
	CreatedAt  time.Time
	Players    map[string]struct{}
}

type Service struct {
	mu           sync.RWMutex
	players      map[string]Player
	nicknameToID map[string]string
	rooms        map[string]*roomRecord
	nextPlayerID int
	nextRoomID   int
}

func NewService() *Service {
	return &Service{
		players:      make(map[string]Player),
		nicknameToID: make(map[string]string),
		rooms:        make(map[string]*roomRecord),
		nextPlayerID: 0,
		nextRoomID:   0,
	}
}

func (s *Service) CreatePlayer(nickname string) (Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cleanNickname := normalizeDisplayNickname(nickname)
	if cleanNickname == "" {
		return Player{}, ErrNicknameRequired
	}

	normalized := normalizeNickname(cleanNickname)
	if existingID, exists := s.nicknameToID[normalized]; exists {
		existingPlayer, ok := s.players[existingID]
		if ok {
			return existingPlayer, nil
		}
	}

	s.nextPlayerID++
	id := fmt.Sprintf("player_%d", s.nextPlayerID)
	player := Player{
		ID:        id,
		Nickname:  cleanNickname,
		CreatedAt: time.Now().UTC(),
	}

	s.players[id] = player
	s.nicknameToID[normalized] = id
	return player, nil
}

func (s *Service) ListPlayers() []Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]Player, 0, len(s.players))
	for _, p := range s.players {
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players
}

func (s *Service) PlayersByIDs(ids []string) []Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]Player, 0, len(ids))
	for _, id := range ids {
		p, ok := s.players[id]
		if !ok {
			continue
		}
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players
}

func (s *Service) GetPlayer(playerID string) (Player, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[strings.TrimSpace(playerID)]
	if !ok {
		return Player{}, false
	}
	return player, true
}

func (s *Service) GetRoom(roomID string) (Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.rooms[strings.TrimSpace(roomID)]
	if !ok {
		return Room{}, false
	}

	return roomSnapshot(room), true
}

func (s *Service) ListRoomsDetailed() []Room {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]Room, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, roomSnapshot(room))
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})

	return rooms
}

func (s *Service) ListRoomViews() []RoomView {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]RoomView, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, RoomView{
			ID:           room.ID,
			Name:         room.Name,
			PlayersCount: len(room.Players),
			MaxPlayers:   room.MaxPlayers,
			IsPrivate:    room.IsPrivate,
		})
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})

	return rooms
}

func (s *Service) CreateRoom(name string, hostPlayerID string, maxPlayers int, isPrivate bool) (Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return Room{}, ErrRoomNameRequired
	}

	hostPlayerID = strings.TrimSpace(hostPlayerID)
	if hostPlayerID == "" {
		return Room{}, ErrInvalidPlayerID
	}

	if _, exists := s.players[hostPlayerID]; !exists {
		return Room{}, ErrPlayerNotFound
	}

	if maxPlayers <= 0 {
		maxPlayers = DefaultMaxPlayers
	}

	s.nextRoomID++
	roomID := fmt.Sprintf("room_%d", s.nextRoomID)

	record := &roomRecord{
		ID:         roomID,
		Name:       name,
		HostID:     hostPlayerID,
		Status:     RoomStatusOpen,
		MaxPlayers: maxPlayers,
		IsPrivate:  isPrivate,
		CreatedAt:  time.Now().UTC(),
		Players: map[string]struct{}{
			hostPlayerID: {},
		},
	}

	s.rooms[roomID] = record
	return roomSnapshot(record), nil
}

func (s *Service) DeleteRoom(roomID string, requesterID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ErrInvalidRoomID
	}
	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" {
		return ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}
	if room.HostID != requesterID {
		return ErrNotRoomHost
	}

	delete(s.rooms, roomID)
	return nil
}

func (s *Service) JoinRoom(roomID string, playerID string) (RoomView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return RoomView{}, ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return RoomView{}, ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return RoomView{}, ErrRoomNotFound
	}
	if room.Status != RoomStatusOpen {
		return RoomView{}, ErrRoomNotOpen
	}
	if _, exists := room.Players[playerID]; exists {
		return roomView(room), nil
	}
	if len(room.Players) >= room.MaxPlayers {
		return RoomView{}, ErrRoomFull
	}

	if _, exists := s.players[playerID]; !exists {
		s.players[playerID] = Player{
			ID:        playerID,
			Nickname:  playerID,
			CreatedAt: time.Now().UTC(),
		}
	}

	room.Players[playerID] = struct{}{}
	return roomView(room), nil
}

func (s *Service) LeaveRoom(roomID string, playerID string) (RoomView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return RoomView{}, ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return RoomView{}, ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return RoomView{}, ErrRoomNotFound
	}
	if room.Status != RoomStatusOpen {
		return RoomView{}, ErrRoomNotOpen
	}

	delete(room.Players, playerID)
	if room.HostID == playerID {
		room.HostID = firstPlayerID(room.Players)
	}
	if len(room.Players) == 0 {
		room.Status = RoomStatusClosed
	}

	return roomView(room), nil
}

func roomSnapshot(room *roomRecord) Room {
	playerIDs := make([]string, 0, len(room.Players))
	for playerID := range room.Players {
		playerIDs = append(playerIDs, playerID)
	}
	sort.Strings(playerIDs)

	return Room{
		ID:         room.ID,
		Name:       room.Name,
		HostID:     room.HostID,
		Status:     room.Status,
		MaxPlayers: room.MaxPlayers,
		IsPrivate:  room.IsPrivate,
		PlayerIDs:  playerIDs,
		CreatedAt:  room.CreatedAt,
	}
}

func roomView(room *roomRecord) RoomView {
	return RoomView{
		ID:           room.ID,
		Name:         room.Name,
		PlayersCount: len(room.Players),
		MaxPlayers:   room.MaxPlayers,
		IsPrivate:    room.IsPrivate,
	}
}

func firstPlayerID(players map[string]struct{}) string {
	ids := make([]string, 0, len(players))
	for playerID := range players {
		ids = append(ids, playerID)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}

func normalizeDisplayNickname(raw string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
}

func normalizeNickname(raw string) string {
	return strings.ToLower(normalizeDisplayNickname(raw))
}
