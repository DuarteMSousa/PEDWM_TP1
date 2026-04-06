package repositories

import (
	"backend/internal/domain/player"
	"backend/internal/domain/room"
	"context"
	"log/slog"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomPostgresRepository implements RoomRepository with PostgreSQL.
type RoomPostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRoomPostgresRepository creates a new room repository.
func NewRoomPostgresRepository(pool *pgxpool.Pool) *RoomPostgresRepository {
	return &RoomPostgresRepository{pool: pool}
}

// Save persists or updates a room and its players (transactional upsert).
func (r *RoomPostgresRepository) Save(rm *room.Room) error {
	ctx := context.Background()
	slog.Debug("persisting room", "roomID", rm.ID, "status", rm.Status, "players", len(rm.Players))

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO rooms (id, host_id, status, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET host_id = $2,
		    status = $3,
		    updated_at = NOW()
	`,
		rm.ID,
		rm.HostID,
		string(rm.Status),
		rm.CreatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM room_players WHERE room_id = $1`, rm.ID)
	if err != nil {
		return err
	}

	orderedPlayers := make([]*player.Player, 0, len(rm.Players))
	for _, p := range rm.Players {
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
			INSERT INTO room_players (room_id, user_id, sequence)
			VALUES ($1, $2, $3)
		`,
			rm.ID,
			p.ID,
			sequence,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindByID finds a room by its ID, including the players.
func (r *RoomPostgresRepository) FindByID(id string) (*room.Room, error) {
	ctx := context.Background()

	var (
		rm     room.Room
		status string
	)

	err := r.pool.QueryRow(ctx, `
		SELECT id, host_id, status, created_at
		FROM rooms WHERE id = $1
	`, id).Scan(&rm.ID, &rm.HostID, &status, &rm.CreatedAt)
	if err != nil {
		return nil, err
	}

	rm.Status = room.RoomStatus(status)

	rows, err := r.pool.Query(ctx, `
		SELECT rp.user_id, u.username, rp.sequence
		FROM room_players rp
		JOIN users u ON u.id = rp.user_id
		WHERE rp.room_id = $1
		ORDER BY rp.sequence ASC, rp.created_at ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rm.Players = make(map[string]*player.Player)

	for rows.Next() {
		var (
			id       string
			username string
			sequence int
		)
		if err := rows.Scan(&id, &username, &sequence); err != nil {
			return nil, err
		}
		rm.Players[id] = player.NewPlayer(id, username, sequence)
	}

	return &rm, nil
}

// FindAll returns all rooms.
func (r *RoomPostgresRepository) FindAll() ([]*room.Room, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `SELECT id FROM rooms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*room.Room

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		rm, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}

		result = append(result, rm)
	}

	return result, nil
}

// Delete removes a room by its ID.
func (r *RoomPostgresRepository) Delete(id string) error {
	ctx := context.Background()
	slog.Debug("deleting room from database", "roomID", id)
	_, err := r.pool.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	if err != nil {
		slog.Error("error deleting room", "roomID", id, "error", err)
	}
	return err
}
