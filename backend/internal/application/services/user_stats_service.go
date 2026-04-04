package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/user"
	"errors"
)

var (
	ErrUserStatsNotFound = errors.New("user stats not found")
)

type UserStatsService struct {
	repo     interfaces.UserStatsRepository
	userRepo interfaces.UserRepository
}

func NewUserStatsService(repo interfaces.UserStatsRepository, userRepo interfaces.UserRepository) *UserStatsService {
	return &UserStatsService{repo: repo, userRepo: userRepo}
}

func (s *UserStatsService) GetByUserID(userID string) (*user.UserStats, error) {
	stats, err := s.repo.FindByUserID(userID)
	if err != nil || stats == nil {
		return nil, ErrUserStatsNotFound
	}
	return stats, nil
}

func (s *UserStatsService) RecordGame(userID string, won bool) (*user.UserStats, error) {
	stats, err := s.repo.FindByUserID(userID)
	if err != nil || stats == nil {
		stats = user.NewUserStats(userID)
		if err := s.repo.Save(stats); err != nil {
			return nil, err
		}
	}

	stats.RecordGame(won)

	if err := s.repo.Save(stats); err != nil {
		return nil, err
	}
	return stats, nil
}
