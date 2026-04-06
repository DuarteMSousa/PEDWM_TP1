package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/card"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	events_infrastructure "backend/internal/infrastructure/events"
	"backend/internal/infrastructure/websocket"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrGameNotFound = errors.New("game not found")
)

// GameSnapshot contains the current state of a game visible to a player.
type GameSnapshot struct {
	RoomID          string
	GameID          string
	Status          string
	TrumpSuit       string
	CurrentPlayerID string
	MyHand          []card.Card
	TablePlays      []trick.Play
	Teams           map[string]*team.Team
	Scores          map[string]int
}

// RoomService orquestra operações sobre salas: criar, entrar, sair, iniciar jogo.
type RoomService struct {
	repo            interfaces.RoomRepository
	gameRepo        interfaces.GameRepository
	userRepo        interfaces.UserRepository
	eventService    *EventService
	eventDispatcher *events_infrastructure.EventDispatcher
	hub             *websocket.Hub
}

// NewRoomService cria um novo RoomService.
func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, eventService *EventService, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, eventService: eventService, hub: hub}
}

// CreateRoom cria uma nova sala e configura o event bus e observers.
func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	slog.Info("a criar sala", "hostID", hostID)

	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		slog.Warn("criação de sala falhada: host não encontrado", "hostID", hostID)
		return nil, err
	}

	r, err := room.NewRoom(uuid.New().String(), hostID, user.Username)
	if err != nil {
		slog.Error("erro ao criar sala", "hostID", hostID, "error", err)
		return nil, err
	}

	eventBus := events.NewEventBus()
	r.SetEventBus(eventBus)

	wsObserver := websocket.NewWebSocketObserver(s.hub)
	eventBus.Subscribe(wsObserver)

	if s.eventDispatcher == nil {
		s.eventDispatcher = events_infrastructure.GetEventDispatcherInstance()
	}

	persistenceObserver := events_infrastructure.NewEventPersistanceObserver(s.eventService, s.eventDispatcher)
	eventBus.Subscribe(persistenceObserver)

	s.hub.CreateRoomHub(r)

	slog.Info("sala criada com sucesso", "roomID", r.ID, "hostID", hostID)
	return r, s.repo.Save(r)
}

// JoinRoom adiciona um jogador a uma sala existente.
func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	slog.Info("a entrar na sala", "roomID", roomID, "userID", userID)

	room := s.hub.GetRoom(roomID)
	if room == nil {
		slog.Warn("entrada na sala falhada: sala não encontrada", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		slog.Warn("entrada na sala falhada: utilizador não encontrado", "userID", userID)
		return nil, err
	}

	if err := room.AddPlayer(userID, user.Username); err != nil {
		slog.Warn("erro ao adicionar jogador à sala", "roomID", roomID, "userID", userID, "error", err)
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		slog.Error("erro ao persistir sala após entrada", "roomID", roomID, "error", err)
		return nil, err
	}

	slog.Info("jogador entrou na sala", "roomID", roomID, "userID", userID, "totalPlayers", len(room.Players))
	return room, nil
}

// LeaveRoom remove um jogador de uma sala. Fecha a sala se ficar vazia.
func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	slog.Info("a sair da sala", "roomID", roomID, "userID", userID)

	rm := s.hub.GetRoom(roomID)
	if rm == nil {
		slog.Warn("saída da sala falhada: sala não encontrada", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	if err := rm.RemovePlayer(userID); err != nil {
		slog.Error("erro ao remover jogador da sala", "roomID", roomID, "userID", userID, "error", err)
		return nil, err
	}

	if rm.Status == room.CLOSED {
		slog.Info("sala vazia, a eliminar", "roomID", roomID)
		if err := s.repo.Delete(rm.ID); err != nil {
			slog.Error("erro ao eliminar sala", "roomID", roomID, "error", err)
			return nil, err
		}
		return rm, nil
	}

	if err := s.repo.Save(rm); err != nil {
		slog.Error("erro ao persistir sala após saída", "roomID", roomID, "error", err)
		return nil, err
	}

	slog.Info("jogador saiu da sala", "roomID", roomID, "userID", userID, "remainingPlayers", len(rm.Players))
	return rm, nil
}

// DeleteRoom remove uma sala da persistência.
func (s *RoomService) DeleteRoom(roomID string) error {
	slog.Info("a eliminar sala", "roomID", roomID)
	if err := s.repo.Delete(roomID); err != nil {
		slog.Error("erro ao eliminar sala", "roomID", roomID, "error", err)
		return err
	}
	slog.Info("sala eliminada", "roomID", roomID)
	return nil
}

// StartGame inicia um jogo de Sueca na sala indicada.
func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	slog.Info("a iniciar jogo", "roomID", roomID)

	room := s.hub.GetRoom(roomID)
	if room == nil {
		slog.Warn("início de jogo falhado: sala não encontrada", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	if err := room.CreateGame(); err != nil {
		slog.Error("erro ao criar jogo", "roomID", roomID, "error", err)
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		slog.Error("erro ao persistir sala após criação de jogo", "roomID", roomID, "error", err)
		return nil, err
	}
	if room.Game != nil {
		if err := s.gameRepo.Save(room.Game); err != nil {
			slog.Error("erro ao persistir jogo", "roomID", roomID, "gameID", room.Game.ID, "error", err)
			return nil, err
		}
	}

	room.Game.State.Enter()

	slog.Info("jogo iniciado com sucesso", "roomID", roomID, "gameID", room.Game.ID)
	return room, nil
}

// GetRoom devolve uma sala pelo ID (hub em memória ou base de dados).
func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	if r := s.hub.GetRoom(id); r != nil {
		return r, nil
	}
	return s.repo.FindByID(id)
}

// GetRooms devolve todas as salas.
func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}

// GetGameSnapshot devolve o estado atual de um jogo visível para um jogador.
func (s *RoomService) GetGameSnapshot(roomID, playerID string) (*GameSnapshot, error) {
	slog.Debug("a obter snapshot do jogo", "roomID", roomID, "playerID", playerID)
	room := s.hub.GetRoom(roomID)
	if room == nil {
		return nil, ErrRoomNotFound
	}
	if room.Game == nil {
		return nil, ErrGameNotFound
	}

	g := room.Game
	currentRound := g.CurrentRound()
	if currentRound == nil {
		return nil, ErrGameNotFound
	}

	requestedPlayer, err := g.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}

	scores := make(map[string]int, len(g.Score))
	for teamID, value := range g.Score {
		scores[teamID] = value
	}

	myHand := make([]card.Card, 0)
	if requestedPlayer.Hand != nil {
		myHand = append(myHand, requestedPlayer.Hand.Cards...)
	}

	tablePlays := make([]trick.Play, 0)
	currentPlayerID := ""
	if currentRound.CurrentTrick != nil {
		tablePlays = append(tablePlays, currentRound.CurrentTrick.Plays...)
		nextPlayerID, nextErr := currentRound.CurrentTrick.TurnOrder.Next()
		if nextErr == nil {
			currentPlayerID = nextPlayerID
		}
	}

	return &GameSnapshot{
		RoomID:          room.ID,
		GameID:          g.ID.String(),
		Status:          string(g.Status),
		TrumpSuit:       string(currentRound.TrumpSuit),
		CurrentPlayerID: currentPlayerID,
		MyHand:          myHand,
		TablePlays:      tablePlays,
		Scores:          scores,
		Teams:           g.Teams,
	}, nil
}
