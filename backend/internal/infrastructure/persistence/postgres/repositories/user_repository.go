package repositories

import (
	"backend/internal/domain/user"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserPostgresRepository implements UserRepository with PostgreSQL.
type UserPostgresRepository struct {
	pool *pgxpool.Pool
}

// NewUserPostgresRepository creates a new user repository.
func NewUserPostgresRepository(pool *pgxpool.Pool) *UserPostgresRepository {
	return &UserPostgresRepository{pool: pool}
}

// Save persists or updates a user (upsert).
func (r *UserPostgresRepository) Save(u *user.User) error {
	ctx := context.Background()
	slog.Debug("persisting user", "userID", u.ID, "username", u.Username)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, username, password)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET username = $2,
		    password = $3
	`, u.ID, u.Username, u.Password)

	if err != nil {
		slog.Error("error persisting user", "userID", u.ID, "error", err)
	}
	return err
}

// FindByID finds a user by ID.
func (r *UserPostgresRepository) FindByID(id string) (*user.User, error) {
	ctx := context.Background()

	var u user.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		slog.Debug("user not found by ID", "userID", id, "error", err)
		return nil, err
	}

	return &u, nil
}

// FindByUsername finds a user by username.
func (r *UserPostgresRepository) FindByUsername(username string) (*user.User, error) {
	ctx := context.Background()

	var u user.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password
		FROM users WHERE username = $1
	`, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		slog.Debug("user not found by username", "username", username, "error", err)
		return nil, err
	}

	return &u, nil
}
