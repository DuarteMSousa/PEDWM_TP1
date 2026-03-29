package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/user"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repo      interfaces.UserRepository
	statsRepo interfaces.UserStatsRepository
}

func NewUserService(repo interfaces.UserRepository, statsRepo interfaces.UserStatsRepository) *UserService {
	return &UserService{repo: repo, statsRepo: statsRepo}
}

func (s *UserService) Register(username, password string) (*user.User, error) {
	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return nil, ErrUsernameExists
	}

	u, err := user.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(u); err != nil {
		return nil, err
	}

	us := user.NewUserStats(u.ID)

	if err := s.statsRepo.Save(us); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) Login(username, password string) (*user.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil || u == nil {
		return nil, ErrInvalidCredentials
	}

	if !u.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *UserService) GetUser(id string) (*user.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil || u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (s *UserService) GetUserByUsername(username string) (*user.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil || u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}
