package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/infrastructure/websocket"
	"errors"

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

	//Ligar observers ao event bus da sala
	s.hub.CreateRoomHub(r.ID, hostID, user.Username)

	eventBus := events.NewEventBus()
	r.SetEventBus(eventBus)

	wsObserver := websocket.NewWebSocketObserver(s.hub)
	eventBus.Subscribe(wsObserver)

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	if err := r.AddPlayer(userID, user.Username); err != nil {
		return nil, err
	}

	roomHub := s.hub.GetRoomHub(roomID)

	if roomHub == nil {
		return nil, ErrRoomNotFound
	}

	room := roomHub.GetRoom()

	if room == nil {
		return nil, ErrRoomNotFound
	}

	room.AddPlayer(userID, user.Username)

	return r, s.repo.Save(r)
}

func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	if err := r.RemovePlayer(userID); err != nil {
		return nil, err
	}

	roomHub := s.hub.GetRoomHub(roomID)

	if roomHub == nil {
		return nil, ErrRoomNotFound
	}

	room := roomHub.GetRoom()

	if room == nil {
		return nil, ErrRoomNotFound
	}

	room.RemovePlayer(userID)

	return r, s.repo.Save(r)
}

func (s *RoomService) StartGame(roomID string) (*room.Room, error) {

	roomHub := s.hub.GetRoomHub(roomID)

	if roomHub == nil {
		return nil, ErrRoomNotFound
	}

	room := roomHub.GetRoom()

	if room == nil {
		return nil, ErrRoomNotFound
	}

	room.StartGame()

	return room, nil
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}
