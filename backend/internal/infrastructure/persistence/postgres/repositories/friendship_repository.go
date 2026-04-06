package repositories

import (
	"backend/internal/domain/friendship"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FriendshipPostgresRepository implements FriendshipRepository with PostgreSQL.
type FriendshipPostgresRepository struct {
	pool *pgxpool.Pool
}

// NewFriendshipPostgresRepository creates a new friendship repository.
func NewFriendshipPostgresRepository(pool *pgxpool.Pool) *FriendshipPostgresRepository {
	return &FriendshipPostgresRepository{pool: pool}
}

// Save persists or updates a friendship (upsert).
func (r *FriendshipPostgresRepository) Save(f *friendship.Friendship) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO friendships (requester_id, addressee_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (requester_id, addressee_id) DO UPDATE
		SET status     = $3,
		    updated_at = $5
	`, f.RequesterID, f.AddresseeID, string(f.Status), f.CreatedAt, f.UpdatedAt)

	return err
}

// Find finds a friendship by the direction requester → addressee.
func (r *FriendshipPostgresRepository) Find(requesterID, addresseeID string) (*friendship.Friendship, error) {
	ctx := context.Background()

	var (
		f      friendship.Friendship
		status string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT requester_id, addressee_id, status, created_at, updated_at
		FROM friendships
		WHERE requester_id = $1 AND addressee_id = $2
	`, requesterID, addresseeID).Scan(&f.RequesterID, &f.AddresseeID, &status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}

	f.Status = friendship.FriendshipStatus(status)
	return &f, nil
}

// FindByUser returns the accepted friendships of a user.
func (r *FriendshipPostgresRepository) FindByUser(userID string) ([]*friendship.Friendship, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT requester_id, addressee_id, status, created_at, updated_at
		FROM friendships
		WHERE (requester_id = $1 OR addressee_id = $1) AND status = 'ACCEPTED'
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFriendships(rows)
}

// FindPendingForUser returns the pending friendship requests for a user.
func (r *FriendshipPostgresRepository) FindPendingForUser(userID string) ([]*friendship.Friendship, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT requester_id, addressee_id, status, created_at, updated_at
		FROM friendships
		WHERE addressee_id = $1 AND status = 'PENDING'
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFriendships(rows)
}

// Delete removes a friendship.
func (r *FriendshipPostgresRepository) Delete(requesterID, addresseeID string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx, `
		DELETE FROM friendships
		WHERE requester_id = $1 AND addressee_id = $2
	`, requesterID, addresseeID)
	return err
}

// scannable is an interface that abstracts both pgx.Rows and pgx.Row for scanning.
type scannable interface {
	Scan(dest ...any) error
	Next() bool
	Close()
}

// scanFriendships is a helper function to scan multiple friendship rows.
func scanFriendships(rows scannable) ([]*friendship.Friendship, error) {
	var result []*friendship.Friendship

	for rows.Next() {
		var (
			f      friendship.Friendship
			status string
		)
		if err := rows.Scan(&f.RequesterID, &f.AddresseeID, &status, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		f.Status = friendship.FriendshipStatus(status)
		result = append(result, &f)
	}

	return result, nil
}
