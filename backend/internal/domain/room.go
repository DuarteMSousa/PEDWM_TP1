package domain

import (
	"errors"
	"strings"
	"time"
)

type RoomStatus string

const (
	SalaAberta  RoomStatus = "ABERTA"
	SalaEmJogo  RoomStatus = "EM_PARTIDA"
	SalaFechada RoomStatus = "FECHADA"
)

var (
	ErrInvalidRoomID          = errors.New("invalid room id")
	ErrInvalidHost            = errors.New("invalid host")
	ErrRoomNotOpen            = errors.New("cannot join/leave: room is not open")
	ErrRoomFull               = errors.New("room is full")
	ErrPlayerAlreadyInRoom    = errors.New("player already in room")
	ErrPlayerNotFoundInRoom   = errors.New("player not found in room")
	ErrCannotStartGamePlayers = errors.New("cannot start game: need exactly 4 players")
	ErrInvalidGameID          = errors.New("invalid game id")
)

// Room representa uma sala (lobby) para formar um jogo.
type Room struct {
	ID      string
	HostID  string
	Players map[string]*Player
	Status  RoomStatus
	GameID  string

	CreatedAt time.Time
}

// NewRoom cria uma sala nova com um host.
// Invariantes:
//   - roomID não vazio
//   - host não nil e host.ID não vazio
func NewRoom(roomID string, host *Player) (*Room, error) {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil, ErrInvalidRoomID
	}
	if host == nil || strings.TrimSpace(host.ID) == "" {
		return nil, ErrInvalidHost
	}

	players := map[string]*Player{
		host.ID: host,
	}

	return &Room{
		ID:        roomID,
		HostID:    host.ID,
		Players:   players,
		Status:    SalaAberta,
		CreatedAt: time.Now(),
	}, nil
}

// AddPlayer adiciona um jogador à sala.
// Regras:
//   - só é possível entrar se a sala estiver OPEN
//   - máximo de 4 jogadores
//   - não permite duplicados por ID
func (r *Room) AddPlayer(player *Player) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != SalaAberta {
		return ErrRoomNotOpen
	}
	if player == nil || strings.TrimSpace(player.ID) == "" {
		return ErrInvalidPlayerID
	}
	if len(r.Players) >= 4 {
		return ErrRoomFull
	}
	if _, exists := r.Players[player.ID]; exists {
		return ErrPlayerAlreadyInRoom
	}

	r.Players[player.ID] = player
	return nil
}

// RemovePlayer remove um jogador da sala.
// Regras:
//   - só é possível sair se a sala estiver OPEN
//   - se o host sair, outro host é promovido de forma determinística
//   - se a sala ficar vazia, fecha automaticamente
func (r *Room) RemovePlayer(playerID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if r.Status != SalaAberta {
		return ErrRoomNotOpen
	}

	playerID = strings.TrimSpace(playerID)
	if playerID == "" {
		return ErrInvalidPlayerID
	}

	if _, exists := r.Players[playerID]; !exists {
		return ErrPlayerNotFoundInRoom
	}

	delete(r.Players, playerID)

	// Se host sair, promover outro automaticamente (determinístico: menor ID lexicográfico)
	if r.HostID == playerID {
		r.HostID = ""
		for id := range r.Players {
			if r.HostID == "" || id < r.HostID {
				r.HostID = id
			}
		}
	}

	// Se sala ficar vazia, fechar
	if len(r.Players) == 0 {
		r.Status = SalaFechada
	}

	return nil
}

// CanStartGame indica se a sala tem condições para iniciar o jogo.
func (r *Room) CanStartGame() bool {
	if r == nil {
		return false
	}
	return r.Status == SalaAberta && len(r.Players) == 4
}

// StartGame transita a sala para EM_PARTIDA e associa o GameID.
// Pré-condição: exatamente 4 jogadores e sala ABERTA.
func (r *Room) StartGame(gameID string) error {
	if r == nil {
		return errors.New("room is nil")
	}
	if !r.CanStartGame() {
		return ErrCannotStartGamePlayers
	}
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return ErrInvalidGameID
	}

	r.Status = SalaEmJogo
	r.GameID = gameID
	return nil
}

// Close fecha a sala (não falha).
func (r *Room) Close() {
	if r == nil {
		return
	}
	r.Status = SalaFechada
}
