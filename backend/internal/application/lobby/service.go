package lobby

import (
	domainevents "backend/internal/domain/events"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DefaultMaxPlayers = 4
	LobbyRoomID       = "lobby"
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

	ErrInvalidRoomID        = errors.New("room id is required")
	ErrRoomNameRequired     = errors.New("room name is required")
	ErrRoomNotFound         = errors.New("room not found")
	ErrRoomNotOpen          = errors.New("cannot join/leave: room is not open")
	ErrRoomFull             = errors.New("room is full")
	ErrNotRoomHost          = errors.New("only the room host can delete this room")
	ErrRoomPasswordRequired = errors.New("password is required for private room")
	ErrRoomPasswordInvalid  = errors.New("invalid room password")
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
	Password   string
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
	eventBus     *domainevents.EventBus
}

func NewService() *Service {
	return NewServiceWithEventBus(nil)
}

func NewServiceWithEventBus(eventBus *domainevents.EventBus) *Service {
	return &Service{
		players:      make(map[string]Player),
		nicknameToID: make(map[string]string),
		rooms:        make(map[string]*roomRecord),
		nextPlayerID: 0,
		nextRoomID:   0,
		eventBus:     eventBus,
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

func (s *Service) CreateRoom(
	name string,
	hostPlayerID string,
	maxPlayers int,
	isPrivate bool,
	password string,
) (Room, error) {
	s.mu.Lock()

	name = strings.TrimSpace(name)
	if name == "" {
		s.mu.Unlock()
		return Room{}, ErrRoomNameRequired
	}

	hostPlayerID = strings.TrimSpace(hostPlayerID)
	if hostPlayerID == "" {
		s.mu.Unlock()
		return Room{}, ErrInvalidPlayerID
	}

	if _, exists := s.players[hostPlayerID]; !exists {
		s.mu.Unlock()
		return Room{}, ErrPlayerNotFound
	}

	if maxPlayers <= 0 {
		maxPlayers = DefaultMaxPlayers
	}

	password = strings.TrimSpace(password)
	if isPrivate && password == "" {
		s.mu.Unlock()
		return Room{}, ErrRoomPasswordRequired
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
		Password:   password,
		CreatedAt:  time.Now().UTC(),
		Players: map[string]struct{}{
			hostPlayerID: {},
		},
	}

	s.rooms[roomID] = record
	snapshot := roomSnapshot(record)
	s.mu.Unlock()

	s.publish(
		domainevents.NewRoomCreatedEvent(
			LobbyRoomID,
			hostPlayerID,
		).WithPayload(roomPayload(snapshot)),
	)
	s.publish(domainevents.NewRoomCreatedEvent(snapshot.ID, hostPlayerID))

	return snapshot, nil
}

func (s *Service) DeleteRoom(roomID string, requesterID string) error {
	s.mu.Lock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		s.mu.Unlock()
		return ErrInvalidRoomID
	}
	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" {
		s.mu.Unlock()
		return ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		s.mu.Unlock()
		return ErrRoomNotFound
	}
	if room.HostID != requesterID {
		s.mu.Unlock()
		return ErrNotRoomHost
	}

	snapshot := roomSnapshot(room)
	delete(s.rooms, roomID)
	s.mu.Unlock()

	deletedPayload := map[string]any{
		"roomId": roomID,
	}
	s.publish(domainevents.NewRoomDeletedEvent(roomID, deletedPayload))
	s.publish(
		domainevents.NewRoomDeletedEvent(
			LobbyRoomID,
			map[string]any{
				"roomId": roomID,
				"room":   roomPayload(snapshot),
			},
		),
	)
	return nil
}

func (s *Service) JoinRoom(roomID string, playerID string, password string) (RoomView, error) {
	s.mu.Lock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		s.mu.Unlock()
		return RoomView{}, ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		s.mu.Unlock()
		return RoomView{}, ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		s.mu.Unlock()
		return RoomView{}, ErrRoomNotFound
	}
	if room.Status != RoomStatusOpen {
		s.mu.Unlock()
		return RoomView{}, ErrRoomNotOpen
	}
	if _, exists := room.Players[playerID]; exists {
		view := roomView(room)
		s.mu.Unlock()
		return view, nil
	}
	if len(room.Players) >= room.MaxPlayers {
		s.mu.Unlock()
		return RoomView{}, ErrRoomFull
	}
	if room.IsPrivate {
		password = strings.TrimSpace(password)
		if password == "" {
			s.mu.Unlock()
			return RoomView{}, ErrRoomPasswordRequired
		}
		if room.Password != password {
			s.mu.Unlock()
			return RoomView{}, ErrRoomPasswordInvalid
		}
	}

	playerNickname := playerID
	if _, exists := s.players[playerID]; !exists {
		s.players[playerID] = Player{
			ID:        playerID,
			Nickname:  playerID,
			CreatedAt: time.Now().UTC(),
		}
	} else {
		playerNickname = s.players[playerID].Nickname
	}

	room.Players[playerID] = struct{}{}
	view := roomView(room)
	snapshot := roomSnapshot(room)
	s.mu.Unlock()

	s.publish(domainevents.NewPlayerJoinedEvent(roomID, playerID, playerNickname))
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			roomID,
			map[string]any{
				"room": roomPayload(snapshot),
			},
		),
	)
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			LobbyRoomID,
			map[string]any{
				"roomId": roomID,
				"room":   roomPayload(snapshot),
			},
		),
	)

	return view, nil
}

func (s *Service) LeaveRoom(roomID string, playerID string) (RoomView, error) {
	s.mu.Lock()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		s.mu.Unlock()
		return RoomView{}, ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		s.mu.Unlock()
		return RoomView{}, ErrInvalidPlayerID
	}

	room, ok := s.rooms[roomID]
	if !ok {
		s.mu.Unlock()
		return RoomView{}, ErrRoomNotFound
	}
	if room.Status != RoomStatusOpen {
		s.mu.Unlock()
		return RoomView{}, ErrRoomNotOpen
	}

	delete(room.Players, playerID)
	if room.HostID == playerID {
		room.HostID = firstPlayerID(room.Players)
	}
	if len(room.Players) == 0 {
		room.Status = RoomStatusClosed
	}

	view := roomView(room)
	snapshot := roomSnapshot(room)
	s.mu.Unlock()

	s.publish(domainevents.NewPlayerLeftEvent(roomID, playerID))
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			roomID,
			map[string]any{
				"room": roomPayload(snapshot),
			},
		),
	)
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			LobbyRoomID,
			map[string]any{
				"roomId": roomID,
				"room":   roomPayload(snapshot),
			},
		),
	)

	return view, nil
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

func (s *Service) publish(event domainevents.Event) {
	if s == nil || s.eventBus == nil {
		return
	}
	s.eventBus.Publish(event)
}

func roomPayload(room Room) map[string]any {
	return map[string]any{
		"id":           room.ID,
		"name":         room.Name,
		"hostPlayerId": room.HostID,
		"status":       room.Status,
		"maxPlayers":   room.MaxPlayers,
		"isPrivate":    room.IsPrivate,
		"playersCount": len(room.PlayerIDs),
		"playerIds":    room.PlayerIDs,
	}
}
