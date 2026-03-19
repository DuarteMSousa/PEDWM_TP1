package postgres

import (
	"backend/internal/application/ports"
	"context"
	"strings"
)

func (s *LobbyStore) GetRoom(roomID string) (ports.Room, bool) {
	ctx := context.Background()
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ports.Room{}, false
	}

	room, ok, err := s.getRoomSnapshot(ctx, roomID)
	if err != nil {
		return ports.Room{}, false
	}
	return room, ok
}

func (s *LobbyStore) ListRoomsDetailed() []ports.Room {
	ctx := context.Background()

	rows, err := s.pool.Query(
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
		ORDER BY r.id`,
	)
	if err != nil {
		return []ports.Room{}
	}
	defer rows.Close()

	rooms := make([]ports.Room, 0)
	for rows.Next() {
		var room ports.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.HostID,
			&room.Status,
			&room.MaxPlayers,
			&room.IsPrivate,
			&room.CreatedAt,
			&room.PlayerIDs,
		); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	return rooms
}

func (s *LobbyStore) ListRoomViews() []ports.RoomView {
	ctx := context.Background()

	rows, err := s.pool.Query(
		ctx,
		`SELECT
			r.id,
			r.name,
			(SELECT COUNT(*) FROM room_players rp WHERE rp.room_id = r.id) AS players_count,
			r.max_players,
			r.is_private
		FROM rooms r
		ORDER BY r.id`,
	)
	if err != nil {
		return []ports.RoomView{}
	}
	defer rows.Close()

	rooms := make([]ports.RoomView, 0)
	for rows.Next() {
		var room ports.RoomView
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.PlayersCount,
			&room.MaxPlayers,
			&room.IsPrivate,
		); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	return rooms
}
