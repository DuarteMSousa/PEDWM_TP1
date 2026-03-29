package repositories

import (
	"backend/internal/domain/user"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStatsPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserStatsPostgresRepository(pool *pgxpool.Pool) *UserStatsPostgresRepository {
	return &UserStatsPostgresRepository{pool: pool}
}

func (r *UserStatsPostgresRepository) Save(stats *user.UserStats) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_stats (user_id, games, wins, elo)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET games = $2,
		    wins  = $3,
		    elo   = $4
	`, stats.UserId, stats.Games, stats.Wins, stats.Elo)

	return err
}

func (r *UserStatsPostgresRepository) FindByUserID(userID string) (*user.UserStats, error) {
	ctx := context.Background()

	var stats user.UserStats
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, games, wins, elo
		FROM user_stats WHERE user_id = $1
	`, userID).Scan(&stats.UserId, &stats.Games, &stats.Wins, &stats.Elo)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
