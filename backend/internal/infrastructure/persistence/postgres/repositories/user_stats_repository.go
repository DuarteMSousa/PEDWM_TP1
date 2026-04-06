package repositories

import (
	"backend/internal/domain/user"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserStatsPostgresRepository implements UserStatsRepository with PostgreSQL.
type UserStatsPostgresRepository struct {
	pool *pgxpool.Pool
}

// NewUserStatsPostgresRepository creates a new user stats repository.
func NewUserStatsPostgresRepository(pool *pgxpool.Pool) *UserStatsPostgresRepository {
	return &UserStatsPostgresRepository{pool: pool}
}

// Save persists or updates a user's statistics (upsert).
func (r *UserStatsPostgresRepository) Save(stats *user.UserStats) error {
	ctx := context.Background()
	slog.Debug("persisting user statistics", "userID", stats.UserId, "games", stats.Games, "elo", stats.Elo)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_stats (user_id, games, wins, elo)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET games = $2,
		    wins  = $3,
		    elo   = $4
	`, stats.UserId, stats.Games, stats.Wins, stats.Elo)

	if err != nil {
		slog.Error("error persisting user statistics", "userID", stats.UserId, "error", err)
	}
	return err
}

// FindByUserID finds a user's statistics by user ID.
func (r *UserStatsPostgresRepository) FindByUserID(userID string) (*user.UserStats, error) {
	ctx := context.Background()

	var stats user.UserStats
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, games, wins, elo
		FROM user_stats WHERE user_id = $1
	`, userID).Scan(&stats.UserId, &stats.Games, &stats.Wins, &stats.Elo)
	if err != nil {
		slog.Debug("user statistics not found", "userID", userID, "error", err)
		return nil, err
	}

	return &stats, nil
}
