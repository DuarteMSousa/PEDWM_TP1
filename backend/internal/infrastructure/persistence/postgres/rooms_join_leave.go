package postgres

import (
	"backend/internal/application/ports"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *LobbyStore) JoinRoom(
	roomID string,
	playerID string,
	password string,
) (ports.RoomView, ports.Room, ports.Player, error) {
	ctx := context.Background()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrInvalidPlayerID
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		status     string
		maxPlayers int
		isPrivate  bool
		dbPassword string
	)
	err = tx.QueryRow(
		ctx,
		`SELECT status, max_players, is_private, COALESCE(password, '')
		 FROM rooms
		 WHERE id = $1
		 FOR UPDATE`,
		roomID,
	).Scan(&status, &maxPlayers, &isPrivate, &dbPassword)
	if err == pgx.ErrNoRows {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotFound
	}
	if err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}

	if status != string(ports.RoomStatusOpen) {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotOpen
	}

	player, ok := s.getPlayerTx(ctx, tx, playerID)
	if !ok {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrPlayerNotFound
	}

	var alreadyInRoom bool
	if err := tx.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM room_players
			WHERE room_id = $1 AND player_id = $2
		)`,
		roomID,
		playerID,
	).Scan(&alreadyInRoom); err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}
	if alreadyInRoom {
		room, ok, err := s.getRoomSnapshotTx(ctx, tx, roomID)
		if err != nil {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, err
		}
		if !ok {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotFound
		}
		if err := tx.Commit(ctx); err != nil {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, err
		}
		return roomViewFromRoom(room), room, player, nil
	}

	var playersCount int
	if err := tx.QueryRow(
		ctx,
		`SELECT COUNT(*)
		 FROM room_players
		 WHERE room_id = $1`,
		roomID,
	).Scan(&playersCount); err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}
	if playersCount >= maxPlayers {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomFull
	}

	if isPrivate {
		password = strings.TrimSpace(password)
		if password == "" {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomPasswordRequired
		}
		if password != dbPassword {
			return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomPasswordInvalid
		}
	}

	if _, err := tx.Exec(
		ctx,
		`INSERT INTO room_players (room_id, player_id, joined_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (room_id, player_id) DO NOTHING`,
		roomID,
		playerID,
	); err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}

	room, ok, err := s.getRoomSnapshotTx(ctx, tx, roomID)
	if err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}
	if !ok {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, ports.ErrRoomNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return ports.RoomView{}, ports.Room{}, ports.Player{}, err
	}

	return roomViewFromRoom(room), room, player, nil
}

func (s *LobbyStore) LeaveRoom(roomID string, playerID string) (ports.RoomView, ports.Room, error) {
	ctx := context.Background()

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.RoomView{}, ports.Room{}, ports.ErrInvalidRoomID
	}
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ports.RoomView{}, ports.Room{}, ports.ErrInvalidPlayerID
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		hostID string
		status string
	)
	err = tx.QueryRow(
		ctx,
		`SELECT host_id, status
		 FROM rooms
		 WHERE id = $1
		 FOR UPDATE`,
		roomID,
	).Scan(&hostID, &status)
	if err == pgx.ErrNoRows {
		return ports.RoomView{}, ports.Room{}, ports.ErrRoomNotFound
	}
	if err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}
	if status != string(ports.RoomStatusOpen) {
		return ports.RoomView{}, ports.Room{}, ports.ErrRoomNotOpen
	}

	if _, err := tx.Exec(
		ctx,
		`DELETE FROM room_players
		 WHERE room_id = $1 AND player_id = $2`,
		roomID,
		playerID,
	); err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}

	var playersCount int
	if err := tx.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM room_players WHERE room_id = $1`,
		roomID,
	).Scan(&playersCount); err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}

	newHostID := hostID
	newStatus := status
	if hostID == playerID {
		if err := tx.QueryRow(
			ctx,
			`SELECT COALESCE((
				SELECT player_id
				FROM room_players
				WHERE room_id = $1
				ORDER BY player_id
				LIMIT 1
			), '')`,
			roomID,
		).Scan(&newHostID); err != nil {
			return ports.RoomView{}, ports.Room{}, err
		}
	}
	if playersCount == 0 {
		newStatus = string(ports.RoomStatusClosed)
		newHostID = ""
	}

	if _, err := tx.Exec(
		ctx,
		`UPDATE rooms
		 SET host_id = $2, status = $3
		 WHERE id = $1`,
		roomID,
		newHostID,
		newStatus,
	); err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}

	room, ok, err := s.getRoomSnapshotTx(ctx, tx, roomID)
	if err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}
	if !ok {
		return ports.RoomView{}, ports.Room{}, ports.ErrRoomNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return ports.RoomView{}, ports.Room{}, err
	}

	return roomViewFromRoom(room), room, nil
}
