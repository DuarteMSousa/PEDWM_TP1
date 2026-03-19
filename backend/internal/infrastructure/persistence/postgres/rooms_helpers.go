package postgres

import (
	"backend/internal/application/ports"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *LobbyStore) getPlayerTx(ctx context.Context, tx pgx.Tx, playerID string) (ports.Player, bool) {
	var player ports.Player
	err := tx.QueryRow(
		ctx,
		`SELECT id, nickname, created_at
		 FROM players
		 WHERE id = $1`,
		playerID,
	).Scan(&player.ID, &player.Nickname, &player.CreatedAt)
	if err == pgx.ErrNoRows {
		return ports.Player{}, false
	}
	if err != nil {
		return ports.Player{}, false
	}
	return player, true
}

func (s *LobbyStore) getRoomSnapshot(ctx context.Context, roomID string) (ports.Room, bool, error) {
	var room ports.Room
	err := s.pool.QueryRow(
		ctx,
		`SELECT
			r.id,
			r.name,
			r.host_id,
			r.status,
			r.max_players,
			r.is_private,
			r.created_at,
			COALESCE(
				(SELECT array_agg(rp.player_id ORDER BY rp.player_id)
				 FROM room_players rp
				 WHERE rp.room_id = r.id),
				ARRAY[]::TEXT[]
			) AS player_ids
		FROM rooms r
		WHERE r.id = $1`,
		roomID,
	).Scan(
		&room.ID,
		&room.Name,
		&room.HostID,
		&room.Status,
		&room.MaxPlayers,
		&room.IsPrivate,
		&room.CreatedAt,
		&room.PlayerIDs,
	)
	if err == pgx.ErrNoRows {
		return ports.Room{}, false, nil
	}
	if err != nil {
		return ports.Room{}, false, err
	}

	return room, true, nil
}

func (s *LobbyStore) getRoomSnapshotTx(ctx context.Context, tx pgx.Tx, roomID string) (ports.Room, bool, error) {
	var room ports.Room
	err := tx.QueryRow(
		ctx,
		`SELECT
			r.id,
			r.name,
			r.host_id,
			r.status,
			r.max_players,
			r.is_private,
			r.created_at,
			COALESCE(
				(SELECT array_agg(rp.player_id ORDER BY rp.player_id)
				 FROM room_players rp
				 WHERE rp.room_id = r.id),
				ARRAY[]::TEXT[]
			) AS player_ids
		FROM rooms r
		WHERE r.id = $1`,
		roomID,
	).Scan(
		&room.ID,
		&room.Name,
		&room.HostID,
		&room.Status,
		&room.MaxPlayers,
		&room.IsPrivate,
		&room.CreatedAt,
		&room.PlayerIDs,
	)
	if err == pgx.ErrNoRows {
		return ports.Room{}, false, nil
	}
	if err != nil {
		return ports.Room{}, false, err
	}

	return room, true, nil
}

func roomViewFromRoom(room ports.Room) ports.RoomView {
	return ports.RoomView{
		ID:           room.ID,
		Name:         room.Name,
		PlayersCount: len(room.PlayerIDs),
		MaxPlayers:   room.MaxPlayers,
		IsPrivate:    room.IsPrivate,
	}
}

func normalizeDisplayNickname(raw string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
}
