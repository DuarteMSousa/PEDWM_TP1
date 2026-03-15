package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(roomID string, host *domain.Player) (*domain.Room, error) {
	room, err := domain.NewRoom(roomID, host)
	if err != nil {
		return nil, err
	}

	s.repo.Save(room)
	return room, nil
}

func (s *RoomService) GetRoom(id string) (*domain.Room, bool) {
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() []*domain.Room {
	return s.repo.FindAll()
}

func (s *RoomService) AddPlayer(roomID string, player *domain.Player) (*domain.Room, error) {
	room, ok := s.repo.FindByID(roomID)
	if !ok {
		return nil, domain.ErrInvalidRoomID
	}

	err := room.AddPlayer(player)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) RemovePlayer(roomID, playerID string) (*domain.Room, error) {
	room, ok := s.repo.FindByID(roomID)
	if !ok {
		return nil, domain.ErrInvalidRoomID
	}

	err := room.RemovePlayer(playerID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) StartGame(roomID, gameID string) (*domain.Room, error) {
	room, ok := s.repo.FindByID(roomID)
	if !ok {
		return nil, domain.ErrInvalidRoomID
	}

	err := room.StartGame(gameID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) CloseRoom(roomID string) bool {
	room, ok := s.repo.FindByID(roomID)
	if !ok {
		return false
	}

	room.Close()
	return true
}
