package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LobbyStore struct {
	pool *pgxpool.Pool
}

func NewLobbyStore(ctx context.Context, databaseURL string) (*LobbyStore, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("database url is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	store := &LobbyStore{pool: pool}
	if err := store.ensureSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return store, nil
}

func (s *LobbyStore) Close() {
	if s == nil || s.pool == nil {
		return
	}
	s.pool.Close()
}
