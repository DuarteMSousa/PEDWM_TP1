package interfaces

import "backend/internal/domain/user"

type UserStatsRepository interface {
	Save(stats *user.UserStats) error
	FindByUserID(userID string) (*user.UserStats, error)
}
