package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/user"
	"errors"
	"log/slog"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService manages user registration, authentication, and retrieval.
type UserService struct {
	repo      interfaces.UserRepository
	statsRepo interfaces.UserStatsRepository
}

// NewUserService creates a new UserService with the injected repositories.
func NewUserService(repo interfaces.UserRepository, statsRepo interfaces.UserStatsRepository) *UserService {
	return &UserService{repo: repo, statsRepo: statsRepo}
}

// Register creates a new user with hashed password and initial statistics.
func (s *UserService) Register(username, password string) (*user.User, error) {
	slog.Info("registering user", "username", username)

	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		slog.Warn("registration failed: username already exists", "username", username)
		return nil, ErrUsernameExists
	}

	u, err := user.NewUser(username, password)
	if err != nil {
		slog.Error("error creating user", "username", username, "error", err)
		return nil, err
	}

	if err := s.repo.Save(u); err != nil {
		slog.Error("error persisting user", "userID", u.ID, "error", err)
		return nil, err
	}

	us := user.NewUserStats(u.ID)

	if err := s.statsRepo.Save(us); err != nil {
		slog.Error("error creating initial statistics", "userID", u.ID, "error", err)
		return nil, err
	}

	slog.Info("user registered successfully", "userID", u.ID, "username", username)
	return u, nil
}

// Login validates a user's credentials.
func (s *UserService) Login(username, password string) (*user.User, error) {
	slog.Debug("login attempt", "username", username)

	u, err := s.repo.FindByUsername(username)
	if err != nil || u == nil {
		slog.Warn("login failed: user not found", "username", username)
		return nil, ErrInvalidCredentials
	}

	if !u.CheckPassword(password) {
		slog.Warn("login failed: invalid password", "username", username)
		return nil, ErrInvalidCredentials
	}

	slog.Info("login successful", "userID", u.ID, "username", username)
	return u, nil
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(id string) (*user.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil || u == nil {
		slog.Debug("user not found", "userID", id)
		return nil, ErrUserNotFound
	}
	return u, nil
}

// GetUserByUsername returns a user by username.
func (s *UserService) GetUserByUsername(username string) (*user.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil || u == nil {
		slog.Debug("user not found", "username", username)
		return nil, ErrUserNotFound
	}
	return u, nil
}
