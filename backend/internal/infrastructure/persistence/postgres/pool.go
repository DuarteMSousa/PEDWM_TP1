package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a new PostgreSQL connection pool using the provided DATABASE_URL.
func NewPostgresPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return pgxpool.New(ctx, url)
}
