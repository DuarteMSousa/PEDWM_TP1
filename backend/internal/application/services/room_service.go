package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/card"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	events_infrastructure "backend/internal/infrastructure/events"
	"backend/internal/infrastructure/websocket"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrGameNotFound = errors.New("game not found")
	ErrNotRoomHost  = errors.New("only the room host can start a game")
)

type GameSnapshot struct {
	RoomID          string
	GameID          string
	Status          string
	TrumpSuit       string
	CurrentPlayerID string
	MyHand          []card.Card
	TablePlays      []trick.Play
	Teams           map[string]*team.Team
	Scores          map[string]int
}

type RoomService struct {
	repo            interfaces.RoomRepository
	gameRepo        interfaces.GameRepository
	userRepo        interfaces.UserRepository
	eventService    *EventService
	eventDispatcher *events_infrastructure.EventDispatcher
	hub             *websocket.Hub
}

func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, eventService *EventService, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, eventService: eventService, hub: hub}
}

func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		return nil, err
	}

	r, err := room.NewRoom(uuid.New().String(), hostID, user.Username)
	if err != nil {
		return nil, err
	}

	s.prepareRoomRuntime(r)

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	rm, err := s.getOrLoadRoom(roomID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	if err := rm.AddPlayer(userID, user.Username); err != nil {
		return nil, err
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(rm)

	return rm, nil
}

func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	rm, err := s.getOrLoadRoom(roomID)
	if err != nil {
		return nil, err
	}

	if err := rm.RemovePlayer(userID); err != nil {
		return nil, err
	}

	if rm.Status == room.CLOSED {
		if err := s.repo.Delete(rm.ID); err != nil {
			return nil, err
		}
		s.hub.DeleteRoomHub(rm.ID)
		return rm, nil
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(rm)

	return rm, nil
}

func (s *RoomService) StartGame(roomID string, requesterID string) (*room.Room, error) {
	rm, err := s.getOrLoadRoom(roomID)
	if err != nil {
		return nil, err
	}

	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" || rm.HostID != requesterID {
		return nil, ErrNotRoomHost
	}
	if _, exists := rm.Players[requesterID]; !exists {
		return nil, ErrNotRoomHost
	}

	if err := rm.StartGame(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	if rm.Game != nil {
		if err := s.gameRepo.Save(rm.Game); err != nil {
			return nil, err
		}
	}
	s.publishRoomSyncedEvent(rm)

	return rm, nil
}

func (s *RoomService) publishRoomSyncedEvent(room *room.Room) {
	if room == nil || room.EventBus == nil {
		return
	}

	room.EventBus.Publish(events.Event{
		ID:        uuid.NewString(),
		Type:      events.EventType("ROOM_SYNCED"),
		RoomID:    room.ID,
		Timestamp: time.Now().UTC(),
		Payload: map[string]any{
			"status":       room.Status,
			"playersCount": len(room.Players),
		},
	})
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	return s.getOrLoadRoom(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	rooms, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, rm := range rooms {
		s.prepareRoomRuntime(rm)
		if rm.SyncStatusWithGame() {
			_ = s.repo.Save(rm)
		}
	}

	return rooms, nil
}

func (s *RoomService) GetGameSnapshot(roomID, playerID string) (*GameSnapshot, error) {
	rm, err := s.getOrLoadRoom(roomID)
	if err != nil {
		return nil, err
	}
	if rm.Game == nil {
		return nil, ErrGameNotFound
	}

	g := rm.Game
	currentRound := g.CurrentRound()
	if currentRound == nil {
		return nil, ErrGameNotFound
	}

	requestedPlayer, err := g.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}

	scores := make(map[string]int, len(g.Score))
	for teamID, value := range g.Score {
		scores[teamID] = value
	}

	myHand := make([]card.Card, 0)
	if requestedPlayer.Hand != nil {
		myHand = append(myHand, requestedPlayer.Hand.Cards...)
	}

	tablePlays := make([]trick.Play, 0)
	currentPlayerID := ""
	if currentRound.CurrentTrick != nil {
		tablePlays = append(tablePlays, currentRound.CurrentTrick.Plays...)
		nextPlayerID, nextErr := currentRound.CurrentTrick.TurnOrder.Next()
		if nextErr == nil {
			currentPlayerID = nextPlayerID
		}
	}

	return &GameSnapshot{
		RoomID:          rm.ID,
		GameID:          g.ID.String(),
		Status:          string(g.Status),
		TrumpSuit:       string(currentRound.TrumpSuit),
		CurrentPlayerID: currentPlayerID,
		MyHand:          myHand,
		TablePlays:      tablePlays,
		Scores:          scores,
		Teams:           g.Teams,
	}, nil
}

func (s *RoomService) HandleClientDisconnect(roomID, userID string) {
	roomID = strings.TrimSpace(roomID)
	userID = strings.TrimSpace(userID)
	if roomID == "" || userID == "" || roomID == "lobby" {
		return
	}

	rm, err := s.getOrLoadRoom(roomID)
	if err != nil || rm == nil {
		return
	}

	if rm.Status == room.OPEN {
		if _, exists := rm.Players[userID]; !exists {
			return
		}
		_, _ = s.LeaveRoom(roomID, userID)
		return
	}

	if rm.SyncStatusWithGame() {
		_ = s.repo.Save(rm)
		if _, exists := rm.Players[userID]; exists {
			_, _ = s.LeaveRoom(roomID, userID)
		}
	}
}

func (s *RoomService) getOrLoadRoom(roomID string) (*room.Room, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil, ErrRoomNotFound
	}

	rm := s.hub.GetRoom(roomID)
	if rm == nil {
		loadedRoom, err := s.repo.FindByID(roomID)
		if err != nil || loadedRoom == nil {
			return nil, ErrRoomNotFound
		}
		rm = loadedRoom
	}

	s.prepareRoomRuntime(rm)

	if rm.SyncStatusWithGame() {
		if err := s.repo.Save(rm); err != nil {
			return nil, err
		}
	}

	return rm, nil
}

func (s *RoomService) prepareRoomRuntime(rm *room.Room) {
	if rm == nil {
		return
	}

	if s.eventDispatcher == nil {
		s.eventDispatcher = events_infrastructure.GetEventDispatcherInstance()
	}

	if rm.EventBus == nil {
		eventBus := events.NewEventBus()
		rm.SetEventBus(eventBus)

		wsObserver := websocket.NewWebSocketObserver(s.hub)
		eventBus.Subscribe(wsObserver)

		persistenceObserver := events_infrastructure.NewEventPersistanceObserver(s.eventService, s.eventDispatcher)
		eventBus.Subscribe(persistenceObserver)
	}

	s.hub.CreateRoomHub(rm)
}
