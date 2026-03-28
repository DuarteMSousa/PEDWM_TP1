package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return pgxpool.New(ctx, url)
}
