package repositories

import (
	"backend/internal/domain/game"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GamePostgresRepository struct {
	pool *pgxpool.Pool
}

func NewGamePostgresRepository(pool *pgxpool.Pool) *GamePostgresRepository {
	return &GamePostgresRepository{pool: pool}
}

func (r *GamePostgresRepository) Save(g *game.Game) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO games (id, room_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET status     = $3,
		    updated_at = $5
	`,
		g.ID.String(),
		g.RoomID,
		string(g.Status),
		g.CreatedAt,
		time.Now().UTC(),
	)

	return err
}

func (r *GamePostgresRepository) FindByID(id string) (*game.Game, error) {
	ctx := context.Background()

	var (
		gameID    string
		roomID    string
		status    string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.pool.QueryRow(ctx, `
		SELECT id, room_id, status, created_at, updated_at
		FROM games WHERE id = $1
	`, id).Scan(&gameID, &roomID, &status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(gameID)
	if err != nil {
		return nil, err
	}

	g := &game.Game{
		ID:        parsedID,
		RoomID:    roomID,
		Status:    game.GameStatus(status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return g, nil
}

func (r *GamePostgresRepository) FindByRoomID(roomID string) ([]*game.Game, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT id, room_id, status, created_at, updated_at
		FROM games WHERE room_id = $1
		ORDER BY created_at DESC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*game.Game

	for rows.Next() {
		var (
			gameID    string
			rID       string
			status    string
			createdAt time.Time
			updatedAt time.Time
		)

		if err := rows.Scan(&gameID, &rID, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		parsedID, err := uuid.Parse(gameID)
		if err != nil {
			return nil, err
		}

		g := &game.Game{
			ID:        parsedID,
			RoomID:    rID,
			Status:    game.GameStatus(status),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		result = append(result, g)
	}

	return result, nil
}
