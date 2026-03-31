package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/infrastructure/websocket"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
)

type RoomService struct {
	repo     interfaces.RoomRepository
	gameRepo interfaces.GameRepository
	userRepo interfaces.UserRepository
	hub      *websocket.Hub
}

func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, hub: hub}
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

	eventBus := events.NewEventBus()
	r.SetEventBus(eventBus)

	wsObserver := websocket.NewWebSocketObserver(s.hub)
	eventBus.Subscribe(wsObserver)

	s.hub.CreateRoomHub(r)

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	room, err := s.ensureRealtimeRoom(roomID)
	if err != nil || room == nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	if err := room.AddPlayer(userID, user.Username); err != nil {
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(room)

	return room, nil
}

func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	room, err := s.ensureRealtimeRoom(roomID)
	if err != nil || room == nil {
		return nil, err
	}

	if err := room.RemovePlayer(userID); err != nil {
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(room)

	return room, nil
}

func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	room, err := s.ensureRealtimeRoom(roomID)
	if err != nil || room == nil {
		return nil, ErrRoomNotFound
	}

	if err := room.StartGame(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		return nil, err
	}
	if room.Game != nil {
		if err := s.gameRepo.Save(room.Game); err != nil {
			return nil, err
		}
	}
	s.publishRoomSyncedEvent(room)

	return room, nil
}

func (s *RoomService) ensureRealtimeRoom(roomID string) (*room.Room, error) {
	roomHub := s.hub.GetRoomHub(roomID)
	if roomHub != nil {
		room := roomHub.GetRoom()
		if room != nil {
			if room.EventBus == nil {
				eventBus := events.NewEventBus()
				eventBus.Subscribe(websocket.NewWebSocketObserver(s.hub))
				room.SetEventBus(eventBus)
			}
			return room, nil
		}
	}

	room, err := s.repo.FindByID(roomID)
	if err != nil || room == nil {
		return nil, ErrRoomNotFound
	}

	eventBus := events.NewEventBus()
	eventBus.Subscribe(websocket.NewWebSocketObserver(s.hub))
	room.SetEventBus(eventBus)

	s.hub.CreateRoomHub(room)
	return room, nil
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
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}
