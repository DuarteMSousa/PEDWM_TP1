package application

import (
	"backend/internal/application/interfaces"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/room"
)

type RoomService struct {
	repo     interfaces.RoomRepository
	gameRepo interfaces.GameRepository
	userRepo interfaces.UserRepository
}

func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo}
}

func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		return nil, err
	}

	r, err := room.NewRoom(hostID, user.Username)
	if err != nil {
		return nil, err
	}

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

	return r, s.repo.Save(r)
}

func (s *RoomService) StartGame(roomID string, botStrategy bot_strategy.IBotStrategy) (*room.Room, error) {
	r, err := s.repo.FindByID(roomID)
	if err != nil || r == nil {
		return nil, err
	}

	if err := r.StartGame(botStrategy); err != nil {
		return nil, err
	}

	if err := s.repo.Save(r); err != nil {
		return nil, err
	}

	if r.Game != nil {
		if err := s.gameRepo.Save(r.Game); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}
