package postgres

import (
	"backend/internal/application/ports"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *LobbyStore) CreateRoom(
	name string,
	hostPlayerID string,
	maxPlayers int,
	isPrivate bool,
	password string,
) (ports.Room, error) {
	ctx := context.Background()

	name = strings.TrimSpace(name)
	if name == "" {
		return ports.Room{}, ports.ErrRoomNameRequired
	}

	hostPlayerID = strings.TrimSpace(hostPlayerID)
	if hostPlayerID == "" {
		return ports.Room{}, ports.ErrInvalidPlayerID
	}

	host, exists := s.GetPlayer(hostPlayerID)
	if !exists || host.ID == "" {
		return ports.Room{}, ports.ErrPlayerNotFound
	}

	if maxPlayers <= 0 {
		maxPlayers = ports.DefaultMaxPlayers
	}

	password = strings.TrimSpace(password)
	if isPrivate && password == "" {
		return ports.Room{}, ports.ErrRoomPasswordRequired
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ports.Room{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var room ports.Room
	err = tx.QueryRow(
		ctx,
		`INSERT INTO rooms (id, name, host_id, status, max_players, is_private, password, created_at)
		 VALUES (concat('room_', nextval('room_seq')), $1, $2, $3, $4, $5, $6, NOW())
		 RETURNING id, name, host_id, status, max_players, is_private, created_at`,
		name,
		hostPlayerID,
		string(ports.RoomStatusOpen),
		maxPlayers,
		isPrivate,
		password,
	).Scan(
		&room.ID,
		&room.Name,
		&room.HostID,
		&room.Status,
		&room.MaxPlayers,
		&room.IsPrivate,
		&room.CreatedAt,
	)
	if err != nil {
		return ports.Room{}, err
	}

	if _, err := tx.Exec(
		ctx,
		`INSERT INTO room_players (room_id, player_id, joined_at)
		 VALUES ($1, $2, NOW())`,
		room.ID,
		hostPlayerID,
	); err != nil {
		return ports.Room{}, err
	}

	room.PlayerIDs = []string{hostPlayerID}

	if err := tx.Commit(ctx); err != nil {
		return ports.Room{}, err
	}

	return room, nil
}

func (s *LobbyStore) DeleteRoom(roomID string, requesterID string) (ports.Room, error) {
	ctx := context.Background()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.Room{}, ports.ErrInvalidRoomID
	}
	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" {
		return ports.Room{}, ports.ErrInvalidPlayerID
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ports.Room{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var hostID string
	err = tx.QueryRow(
		ctx,
		`SELECT host_id FROM rooms WHERE id = $1 FOR UPDATE`,
		roomID,
	).Scan(&hostID)
	if err == pgx.ErrNoRows {
		return ports.Room{}, ports.ErrRoomNotFound
	}
	if err != nil {
		return ports.Room{}, err
	}

	if hostID != requesterID {
		return ports.Room{}, ports.ErrNotRoomHost
	}

	snapshot, ok, err := s.getRoomSnapshotTx(ctx, tx, roomID)
	if err != nil {
		return ports.Room{}, err
	}
	if !ok {
		return ports.Room{}, ports.ErrRoomNotFound
	}

	if _, err := tx.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, roomID); err != nil {
		return ports.Room{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return ports.Room{}, err
	}

	return snapshot, nil
}
