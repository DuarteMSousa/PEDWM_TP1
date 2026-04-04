package repositories

import (
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	"context"
	"sort"
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

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
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

	gamePlayers := []*player.Player{}

	for _, team := range g.Teams {
		for _, player := range team.Players {
			gamePlayers = append(gamePlayers, player)
		}
	}

	orderedPlayers := make([]*player.Player, 0, len(gamePlayers))
	for _, p := range gamePlayers {
		orderedPlayers = append(orderedPlayers, p)
	}
	sort.Slice(orderedPlayers, func(i, j int) bool {
		if orderedPlayers[i].Sequence == orderedPlayers[j].Sequence {
			return orderedPlayers[i].ID < orderedPlayers[j].ID
		}
		return orderedPlayers[i].Sequence < orderedPlayers[j].Sequence
	})

	for idx, p := range orderedPlayers {
		sequence := p.Sequence
		if sequence <= 0 {
			sequence = idx + 1
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO game_players (game_id, user_id, sequence)
			VALUES ($1, $2, $3)
		`,
			g.ID,
			p.ID,
			sequence,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
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

func (r *GamePostgresRepository) GetByUserID(userID string) ([]*game.Game, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.room_id, g.status, g.created_at, g.updated_at
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.user_id = $1
		ORDER BY g.created_at DESC
	`, userID)
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
