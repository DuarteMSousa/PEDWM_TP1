package events_infrastructure

import "backend/internal/domain/user"

type UserStatsService interface {
	RecordGame(userID string, won bool) (*user.UserStats, error)
}
