package repositories

import (
	"backend/internal/domain/user"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserPostgresRepository(pool *pgxpool.Pool) *UserPostgresRepository {
	return &UserPostgresRepository{pool: pool}
}

func (r *UserPostgresRepository) Save(u *user.User) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, username, password)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET username = $2,
		    password = $3
	`, u.ID, u.Username, u.Password)

	return err
}

func (r *UserPostgresRepository) FindByID(id string) (*user.User, error) {
	ctx := context.Background()

	var u user.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgresRepository) FindByUsername(username string) (*user.User, error) {
	ctx := context.Background()

	var u user.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, password
		FROM users WHERE username = $1
	`, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
