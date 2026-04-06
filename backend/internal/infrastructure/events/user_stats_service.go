package events_infrastructure

import "backend/internal/domain/user"

// UserStatsService defines the interface for interacting with the user stats service.
type UserStatsService interface {
	RecordGame(userID string, won bool) (*user.UserStats, error)
}
