package postgres

import (
	"backend/internal/application/ports"
	"context"
	"sort"
	"strings"
)

func (s *LobbyStore) CreateOrGetByNickname(nickname string) (ports.Player, error) {
	ctx := context.Background()

	cleanNickname := normalizeDisplayNickname(nickname)
	if cleanNickname == "" {
		return ports.Player{}, ports.ErrNicknameRequired
	}

	normalized := strings.ToLower(cleanNickname)

	var player ports.Player
	err := s.pool.QueryRow(
		ctx,
		`INSERT INTO players (id, nickname, nickname_normalized, created_at)
		 VALUES (concat('player_', nextval('player_seq')), $1, $2, NOW())
		 ON CONFLICT (nickname_normalized)
		 DO UPDATE SET nickname = players.nickname
		 RETURNING id, nickname, created_at`,
		cleanNickname,
		normalized,
	).Scan(&player.ID, &player.Nickname, &player.CreatedAt)
	if err != nil {
		return ports.Player{}, err
	}

	return player, nil
}

func (s *LobbyStore) GetPlayer(playerID string) (ports.Player, bool) {
	ctx := context.Background()
	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ports.Player{}, false
	}

	var player ports.Player
	err := s.pool.QueryRow(
		ctx,
		`SELECT id, nickname, created_at FROM players WHERE id = $1`,
		playerID,
	).Scan(&player.ID, &player.Nickname, &player.CreatedAt)
	if err != nil {
		return ports.Player{}, false
	}

	return player, true
}

func (s *LobbyStore) ListPlayers() []ports.Player {
	ctx := context.Background()
	rows, err := s.pool.Query(
		ctx,
		`SELECT id, nickname, created_at FROM players ORDER BY id`,
	)
	if err != nil {
		return []ports.Player{}
	}
	defer rows.Close()

	players := make([]ports.Player, 0)
	for rows.Next() {
		var player ports.Player
		if err := rows.Scan(&player.ID, &player.Nickname, &player.CreatedAt); err != nil {
			continue
		}
		players = append(players, player)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players
}

func (s *LobbyStore) PlayersByIDs(ids []string) []ports.Player {
	ctx := context.Background()
	if len(ids) == 0 {
		return []ports.Player{}
	}

	rows, err := s.pool.Query(
		ctx,
		`SELECT id, nickname, created_at
		 FROM players
		 WHERE id = ANY($1)
		 ORDER BY id`,
		ids,
	)
	if err != nil {
		return []ports.Player{}
	}
	defer rows.Close()

	players := make([]ports.Player, 0, len(ids))
	for rows.Next() {
		var player ports.Player
		if err := rows.Scan(&player.ID, &player.Nickname, &player.CreatedAt); err != nil {
			continue
		}
		players = append(players, player)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players
}
