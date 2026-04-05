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
	"log"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrGameNotFound = errors.New("game not found")
)

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

type RoomService struct {
	repo            interfaces.RoomRepository
	gameRepo        interfaces.GameRepository
	userRepo        interfaces.UserRepository
	eventService    *EventService
	eventDispatcher *events_infrastructure.EventDispatcher
	hub             *websocket.Hub
}

func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, eventService *EventService, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, eventService: eventService, hub: hub}
}

func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		return nil, err
	}

	r, err := room.NewRoom(uuid.New().String(), hostID, user.Username)
	if err != nil {
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

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	room := s.hub.GetRoom(roomID)
	if room == nil {
		return nil, ErrRoomNotFound
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	if err := room.AddPlayer(userID, user.Username); err != nil {
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	rm := s.hub.GetRoom(roomID)
	if rm == nil {
		return nil, ErrRoomNotFound
	}

	if err := rm.RemovePlayer(userID); err != nil {
		log.Printf("error in leaveroom: %v", err)
		return nil, err
	}

	if rm.Status == room.CLOSED {
		if err := s.repo.Delete(rm.ID); err != nil {
			return nil, err
		}
		return rm, nil
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	return rm, nil
}

func (s *RoomService) DeleteRoom(roomID string) error {
	if err := s.repo.Delete(roomID); err != nil {
		return err
	}
	return nil
}

func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	room := s.hub.GetRoom(roomID)
	if room == nil {
		return nil, ErrRoomNotFound
	}

	if err := room.CreateGame(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		return nil, err
	}
	if room.Game != nil {
		if err := s.gameRepo.Save(room.Game); err != nil {
			return nil, err
		}
	}

	room.Game.State.Enter()

	return room, nil
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	if r := s.hub.GetRoom(id); r != nil {
		return r, nil
	}
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}

func (s *RoomService) GetGameSnapshot(roomID, playerID string) (*GameSnapshot, error) {
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
