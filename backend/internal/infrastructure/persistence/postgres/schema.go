package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if err := cleanupLegacyLobbySchema(ctx, pool); err != nil {
		return err
	}

	statements := []string{

		// =========================
		// ENUM TYPES
		// =========================
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'friendship_status') THEN
				CREATE TYPE friendship_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');
			END IF;
		END$$;`,

		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'room_status') THEN
				CREATE TYPE room_status AS ENUM ('OPEN', 'IN_GAME', 'CLOSED');
			END IF;
		END$$;`,

		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'game_status') THEN
				CREATE TYPE game_status AS ENUM ('PENDING', 'IN_PROGRESS', 'FINISHED');
			END IF;
		END$$;`,

		// =========================
		// USERS
		// =========================
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR NOT NULL UNIQUE,
			password VARCHAR NOT NULL
		)`,

		// =========================
		// FRIENDSHIPS
		// =========================
		`CREATE TABLE IF NOT EXISTS friendships (
			requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			addressee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status friendship_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (requester_id, addressee_id)
		)`,

		// =========================
		// USER STATS
		// =========================
		`CREATE TABLE IF NOT EXISTS user_stats (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			games INT NOT NULL DEFAULT 0,
			wins INT NOT NULL DEFAULT 0,
			elo INT NOT NULL DEFAULT 1000
		)`,

		// =========================
		// ROOMS
		// =========================
		`CREATE TABLE IF NOT EXISTS rooms (
			id UUID PRIMARY KEY,
			host_id UUID NOT NULL REFERENCES users(id),
			status room_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// =========================
		// ROOM PLAYERS
		// =========================
		`CREATE TABLE IF NOT EXISTS room_players (
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (room_id, user_id)
		)`,

		// =========================
		// GAMES
		// =========================
		`CREATE TABLE IF NOT EXISTS games (
			id UUID PRIMARY KEY,
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			status game_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// =========================
		// EVENTS
		// =========================
		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE SET NULL,
			event_type TEXT NOT NULL,
			sequence INT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			payload JSONB
		)`,

		// =========================
		// INDEXES
		// =========================
		`CREATE INDEX IF NOT EXISTS idx_friendships_requester ON friendships(requester_id)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_addressee ON friendships(addressee_id)`,

		`CREATE INDEX IF NOT EXISTS idx_rooms_host ON rooms(host_id)`,

		`CREATE INDEX IF NOT EXISTS idx_games_room_id ON games(room_id)`,

		`CREATE INDEX IF NOT EXISTS idx_room_players_room ON room_players(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_room_players_user ON room_players(user_id)`,

		`CREATE INDEX IF NOT EXISTS idx_events_game ON events(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_events_user ON events(user_id)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS one_active_game_per_room
			ON games(room_id)
			WHERE status = 'IN_PROGRESS'`,
	}

	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func cleanupLegacyLobbySchema(ctx context.Context, pool *pgxpool.Pool) error {
	isLegacy, err := hasLegacyLobbySchema(ctx, pool)
	if err != nil || !isLegacy {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	statements := []string{
		`DROP TABLE IF EXISTS events`,
		`DROP TABLE IF EXISTS games`,
		`DROP TABLE IF EXISTS room_players`,
		`DROP TABLE IF EXISTS rooms`,
		`DROP TYPE IF EXISTS room_status`,
		`DROP TYPE IF EXISTS game_status`,
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func hasLegacyLobbySchema(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var roomsExists bool
	if err := pool.QueryRow(ctx, `SELECT to_regclass('public.rooms') IS NOT NULL`).Scan(&roomsExists); err != nil {
		return false, err
	}
	if !roomsExists {
		return false, nil
	}

	var roomsIDType string
	if err := pool.QueryRow(ctx, `
		SELECT data_type
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = 'rooms'
		  AND column_name = 'id'
	`).Scan(&roomsIDType); err != nil {
		return false, err
	}
	if strings.TrimSpace(strings.ToLower(roomsIDType)) != "uuid" {
		return true, nil
	}

	var hasLegacyRoomColumns bool
	if err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public'
			  AND table_name = 'rooms'
			  AND column_name IN ('name', 'max_players', 'is_private', 'password')
		)
	`).Scan(&hasLegacyRoomColumns); err != nil {
		return false, err
	}
	if hasLegacyRoomColumns {
		return true, nil
	}

	var hasLegacyRoomPlayersColumns bool
	if err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public'
			  AND table_name = 'room_players'
			  AND column_name IN ('player_id', 'joined_at')
		)
	`).Scan(&hasLegacyRoomPlayersColumns); err != nil {
		return false, err
	}
	if hasLegacyRoomPlayersColumns {
		return true, nil
	}

	return false, nil
}
