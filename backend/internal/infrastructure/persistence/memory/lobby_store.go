package memory

import (
	"backend/internal/application/ports"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type roomRecord struct {
	ID         string
	Name       string
	HostID     string
	Status     ports.RoomStatus
	MaxPlayers int
	IsPrivate  bool
	Password   string
	CreatedAt  time.Time
	Players    map[string]struct{}
}

// LobbyStore is an in-memory persistence adapter for lobby player/room data.
type LobbyStore struct {
	mu           sync.RWMutex
	players      map[string]ports.Player
	nicknameToID map[string]string
	rooms        map[string]*roomRecord
	nextPlayerID int
	nextRoomID   int
}

func NewLobbyStore() *LobbyStore {
	return &LobbyStore{
		players:      make(map[string]ports.Player),
		nicknameToID: make(map[string]string),
		rooms:        make(map[string]*roomRecord),
	}
}

func (s *LobbyStore) CreateOrGetByNickname(nickname string) (ports.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cleanNickname := normalizeDisplayNickname(nickname)
	if cleanNickname == "" {
		return ports.Player{}, ports.ErrNicknameRequired
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
	player := ports.Player{
		ID:        id,
		Nickname:  cleanNickname,
		CreatedAt: time.Now().UTC(),
	}

	s.players[id] = player
	s.nicknameToID[normalized] = id
	return player, nil
}

func (s *LobbyStore) GetPlayer(playerID string) (ports.Player, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[strings.TrimSpace(playerID)]
	if !ok {
		return ports.Player{}, false
	}
	return player, true
}

func (s *LobbyStore) ListPlayers() []ports.Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]ports.Player, 0, len(s.players))
	for _, p := range s.players {
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players
}

func (s *LobbyStore) PlayersByIDs(ids []string) []ports.Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]ports.Player, 0, len(ids))
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

func (s *LobbyStore) CreateRoom(
	name string,
	hostPlayerID string,
	maxPlayers int,
	isPrivate bool,
	password string,
) (ports.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return ports.Room{}, ports.ErrRoomNameRequired
	}

	hostPlayerID = strings.TrimSpace(hostPlayerID)
	if hostPlayerID == "" {
		return ports.Room{}, ports.ErrInvalidPlayerID
	}

	if _, exists := s.players[hostPlayerID]; !exists {
		return ports.Room{}, ports.ErrPlayerNotFound
	}

	if maxPlayers <= 0 {
		maxPlayers = ports.DefaultMaxPlayers
	}

	password = strings.TrimSpace(password)
	if isPrivate && password == "" {
		return ports.Room{}, ports.ErrRoomPasswordRequired
	}

	s.nextRoomID++
	roomID := fmt.Sprintf("room_%d", s.nextRoomID)

	record := &roomRecord{
		ID:         roomID,
		Name:       name,
		HostID:     hostPlayerID,
		Status:     ports.RoomStatusOpen,
		MaxPlayers: maxPlayers,
		IsPrivate:  isPrivate,
		Password:   password,
		CreatedAt:  time.Now().UTC(),
		Players: map[string]struct{}{
			hostPlayerID: {},
		},
	}

	s.rooms[roomID] = record
	return roomSnapshot(record), nil
}

func (s *LobbyStore) DeleteRoom(roomID string, requesterID string) (ports.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.Room{}, ports.ErrInvalidRoomID
	}
	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" {
		return ports.Room{}, ports.ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return ports.Room{}, ports.ErrRoomNotFound
	}
	if room.HostID != requesterID {
		return ports.Room{}, ports.ErrNotRoomHost
	}

	snapshot := roomSnapshot(room)
	delete(s.rooms, roomID)
	return snapshot, nil
}

func (s *LobbyStore) JoinRoom(
	roomID string,
	playerID string,
	password string,
) (ports.RoomView, ports.Room, ports.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotFound
	}
	if room.Status != ports.RoomStatusOpen {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotOpen
	}
	if _, exists := room.Players[playerID]; exists {
		player := s.ensurePlayerLocked(playerID)
		return roomView(room), roomSnapshot(room), player, nil
	}
	if len(room.Players) >= room.MaxPlayers {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomFull
	}

	if room.IsPrivate {
		password = strings.TrimSpace(password)
		if password == "" {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomPasswordRequired
		}
		if room.Password != password {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomPasswordInvalid
		}
	}

	player := s.ensurePlayerLocked(playerID)
	room.Players[playerID] = struct{}{}

	return roomView(room), roomSnapshot(room), player, nil
}

func (s *LobbyStore) LeaveRoom(roomID string, playerID string) (ports.RoomView, ports.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.RoomView{}, ports.Room{}, ports.ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ports.RoomView{}, ports.Room{}, ports.ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		return ports.RoomView{}, ports.Room{}, ports.ErrRoomNotFound
	}
	if room.Status != ports.RoomStatusOpen {
		return ports.RoomView{}, ports.Room{}, ports.ErrRoomNotOpen
	}

	delete(room.Players, playerID)
	if room.HostID == playerID {
		room.HostID = firstPlayerID(room.Players)
	}
	if len(room.Players) == 0 {
		room.Status = ports.RoomStatusClosed
	}

	return roomView(room), roomSnapshot(room), nil
}

func (s *LobbyStore) GetRoom(roomID string) (ports.Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.rooms[strings.TrimSpace(roomID)]
	if !ok {
		return ports.Room{}, false
	}
	return roomSnapshot(room), true
}

func (s *LobbyStore) ListRoomsDetailed() []ports.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]ports.Room, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, roomSnapshot(room))
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})

	return rooms
}

func (s *LobbyStore) ListRoomViews() []ports.RoomView {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]ports.RoomView, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, ports.RoomView{
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

func (s *LobbyStore) ensurePlayerLocked(playerID string) ports.Player {
	player, ok := s.players[playerID]
	if ok {
		return player
	}

	player = ports.Player{
		ID:        playerID,
		Nickname:  playerID,
		CreatedAt: time.Now().UTC(),
	}
	s.players[playerID] = player
	return player
}

func roomSnapshot(room *roomRecord) ports.Room {
	playerIDs := make([]string, 0, len(room.Players))
	for playerID := range room.Players {
		playerIDs = append(playerIDs, playerID)
	}
	sort.Strings(playerIDs)

	return ports.Room{
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

func roomView(room *roomRecord) ports.RoomView {
	return ports.RoomView{
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
