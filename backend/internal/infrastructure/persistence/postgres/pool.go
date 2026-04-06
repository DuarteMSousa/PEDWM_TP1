package postgres

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a new PostgreSQL connection pool using the provided DATABASE_URL.
func NewPostgresPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	slog.Info("creating PostgreSQL connection pool")
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		slog.Error("failed to create PostgreSQL connection pool", "error", err)
		return nil, err
	}
	slog.Info("PostgreSQL connection pool created successfully")
	return pool, nil
}
