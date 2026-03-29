package repositories

import (
	"backend/internal/domain/room"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRoomPostgresRepository(pool *pgxpool.Pool) *RoomPostgresRepository {
	return &RoomPostgresRepository{pool: pool}
}

func (r *RoomPostgresRepository) Save(rm *room.Room) error {
	ctx := context.Background()

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

	for _, p := range rm.Players {
		_, err := tx.Exec(ctx, `
			INSERT INTO room_players (room_id, user_id)
			VALUES ($1, $2)
		`,
			rm.ID,
			p.UserID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

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
		SELECT rp.user_id, u.username
		FROM room_players rp
		JOIN users u ON u.id = rp.user_id
		WHERE rp.room_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rm.Players = make(map[string]*room.RoomPlayer)

	for rows.Next() {
		var p room.RoomPlayer
		if err := rows.Scan(&p.UserID, &p.Username); err != nil {
			return nil, err
		}
		rm.Players[p.UserID] = &p
	}

	return &rm, nil
}

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

func (r *RoomPostgresRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	return err
}
