package postgres

import "context"

func (s *LobbyStore) ensureSchema(ctx context.Context) error {
	statements := []string{
		`CREATE SEQUENCE IF NOT EXISTS player_seq START WITH 1 INCREMENT BY 1`,
		`CREATE SEQUENCE IF NOT EXISTS room_seq START WITH 1 INCREMENT BY 1`,
		`CREATE TABLE IF NOT EXISTS players (
			id TEXT PRIMARY KEY,
			nickname TEXT NOT NULL,
			nickname_normalized TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS rooms (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			host_id TEXT NOT NULL REFERENCES players(id),
			status TEXT NOT NULL CHECK (status IN ('OPEN', 'IN_GAME', 'CLOSED')),
			max_players INT NOT NULL,
			is_private BOOLEAN NOT NULL DEFAULT FALSE,
			password TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS room_players (
			room_id TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
			joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (room_id, player_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_room_players_room_id ON room_players(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_room_players_player_id ON room_players(player_id)`,
	}

	for _, stmt := range statements {
		if _, err := s.pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}
