package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/user"
	"errors"
	"log/slog"
)

var (
	ErrUserStatsNotFound = errors.New("user stats not found")
)

// UserStatsService manages user game statistics.
type UserStatsService struct {
	repo     interfaces.UserStatsRepository
	userRepo interfaces.UserRepository
}

// NewUserStatsService creates a new UserStatsService.
func NewUserStatsService(repo interfaces.UserStatsRepository, userRepo interfaces.UserRepository) *UserStatsService {
	return &UserStatsService{repo: repo, userRepo: userRepo}
}

// GetByUserID returns the statistics of a user by their ID.
func (s *UserStatsService) GetByUserID(userID string) (*user.UserStats, error) {
	stats, err := s.repo.FindByUserID(userID)
	if err != nil || stats == nil {
		slog.Debug("user stats not found", "userID", userID)
		return nil, ErrUserStatsNotFound
	}
	return stats, nil
}

// RecordGame records the result of a game for a user,
// updating counts and ELO.
func (s *UserStatsService) RecordGame(userID string, won bool) (*user.UserStats, error) {
	slog.Info("recording game result", "userID", userID, "won", won)

	stats, err := s.repo.FindByUserID(userID)
	if err != nil || stats == nil {
		slog.Debug("user stats not found, creating new", "userID", userID)
		stats = user.NewUserStats(userID)
		if err := s.repo.Save(stats); err != nil {
			slog.Error("error creating user stats", "userID", userID, "error", err)
			return nil, err
		}
	}

	stats.RecordGame(won)

	if err := s.repo.Save(stats); err != nil {
		slog.Error("error persisting user stats", "userID", userID, "error", err)
		return nil, err
	}

	slog.Info("game result recorded", "userID", userID, "games", stats.Games, "wins", stats.Wins, "elo", stats.Elo)
	return stats, nil
}
