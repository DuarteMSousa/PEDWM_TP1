package interfaces

import "backend/internal/domain/user"

// UserStatsRepository defines the contract for user statistics persistence.
type UserStatsRepository interface {
	Save(stats *user.UserStats) error
	FindByUserID(userID string) (*user.UserStats, error)
}
