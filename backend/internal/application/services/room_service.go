package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/player"
	"backend/internal/domain/room"
)

type RoomService struct {
	repo interfaces.RoomRepository
}

func NewRoomService(repo interfaces.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(roomID, hostID, hostName string) (*room.Room, error) {
	host := player.NewPlayer(hostID, hostName, 0)

	r, err := room.NewRoom(roomID, &host)
	if err != nil {
		return nil, err
	}

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, playerID, playerName string) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	p := player.NewPlayer(playerID, playerName, len(r.Players))

	if err := r.AddPlayer(&p); err != nil {
		return nil, err
	}

	return r, s.repo.Save(r)
}

func (s *RoomService) LeaveRoom(roomID, playerID string) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	if err := r.RemovePlayer(playerID); err != nil {
		return nil, err
	}

	return r, s.repo.Save(r)
}

func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	if err := r.StartGame(); err != nil {
		return nil, err
	}

	return r, s.repo.Save(r)
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}
