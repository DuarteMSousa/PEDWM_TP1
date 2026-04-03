package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/events"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

type EventService struct {
	repo interfaces.EventRepository
}

func NewEventService(repo interfaces.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) SaveEvent(event events.Event) error {
	return s.repo.Save(event)
}

type MatchHistoryEntry struct {
	GameID        string
	RoomID        string
	PlayedAt      time.Time
	Won           bool
	WinnerTeamID  string
	MyTeamID      string
	MyScore       int
	OpponentScore int
	FinalScores   map[string]int
}

type ReplayEvent struct {
	ID        string
	Type      string
	Sequence  int
	Timestamp time.Time
	Payload   string
}

type GameReplay struct {
	GameID string
	RoomID string
	Events []ReplayEvent
}

func (s *EventService) GetHistoryByUser(userID string) ([]MatchHistoryEntry, error) {
	allEndedEvents, err := s.repo.FindByType(events.EventGameEnded)
	if err != nil {
		return nil, err
	}

	history := make([]MatchHistoryEntry, 0, len(allEndedEvents))
	for _, endedEvent := range allEndedEvents {
		entry, ok := buildHistoryEntry(endedEvent, userID)
		if !ok {
			continue
		}
		history = append(history, entry)
	}

	return history, nil
}

func (s *EventService) GetReplayByGameID(gameID string) (*GameReplay, error) {
	eventsByGame, err := s.repo.FindByGameID(gameID)
	if err != nil {
		return nil, err
	}
	if len(eventsByGame) == 0 {
		return nil, nil
	}

	replayEvents := make([]ReplayEvent, 0, len(eventsByGame))
	for _, e := range eventsByGame {
		payloadJSON := "{}"
		if e.Payload != nil {
			b, marshalErr := json.Marshal(e.Payload)
			if marshalErr != nil {
				return nil, fmt.Errorf("marshal replay payload: %w", marshalErr)
			}
			payloadJSON = string(b)
		}

		replayEvents = append(replayEvents, ReplayEvent{
			ID:        e.ID,
			Type:      string(e.Type),
			Sequence:  e.Sequence,
			Timestamp: e.Timestamp,
			Payload:   payloadJSON,
		})
	}

	return &GameReplay{
		GameID: gameID,
		RoomID: eventsByGame[0].RoomID,
		Events: replayEvents,
	}, nil
}

func buildHistoryEntry(endedEvent events.Event, userID string) (MatchHistoryEntry, bool) {
	payloadMap, ok := toStringAnyMap(endedEvent.Payload)
	if !ok {
		return MatchHistoryEntry{}, false
	}

	finalScores := parseScoreMap(payloadMap["finalScores"])
	if len(finalScores) == 0 {
		finalScores = parseScoreMap(payloadMap["score"])
	}

	teamIDs := make([]string, 0, len(finalScores))
	for teamID := range finalScores {
		teamIDs = append(teamIDs, teamID)
	}
	slices.Sort(teamIDs)

	myTeamID := findUserTeamID(payloadMap["teams"], userID)
	if myTeamID == "" {
		return MatchHistoryEntry{}, false
	}

	opponentScore := 0
	for _, teamID := range teamIDs {
		if teamID == myTeamID {
			continue
		}
		opponentScore = finalScores[teamID]
		break
	}

	winnerTeamID := toString(payloadMap["winner"])
	return MatchHistoryEntry{
		GameID:        endedEvent.GameID,
		RoomID:        endedEvent.RoomID,
		PlayedAt:      endedEvent.Timestamp,
		WinnerTeamID:  winnerTeamID,
		MyTeamID:      myTeamID,
		Won:           strings.EqualFold(winnerTeamID, myTeamID),
		MyScore:       finalScores[myTeamID],
		OpponentScore: opponentScore,
		FinalScores:   finalScores,
	}, true
}

func findUserTeamID(teamsRaw any, userID string) string {
	teamsMap, ok := toStringAnyMap(teamsRaw)
	if !ok {
		return ""
	}

	for teamID, teamRaw := range teamsMap {
		teamMap, teamOK := toStringAnyMap(teamRaw)
		if !teamOK {
			continue
		}

		playersRaw, hasPlayers := mapGetCaseInsensitive(teamMap, "players")
		if !hasPlayers {
			continue
		}

		players, okPlayers := playersRaw.([]any)
		if !okPlayers {
			continue
		}

		for _, rawPlayer := range players {
			playerMap, playerOK := toStringAnyMap(rawPlayer)
			if !playerOK {
				continue
			}

			playerID, _ := mapGetCaseInsensitive(playerMap, "id")
			if toString(playerID) == userID {
				return teamID
			}
		}
	}

	return ""
}

func parseScoreMap(raw any) map[string]int {
	scoreMap, ok := toStringAnyMap(raw)
	if !ok {
		return map[string]int{}
	}

	result := make(map[string]int, len(scoreMap))
	for teamID, scoreRaw := range scoreMap {
		scoreValue, ok := toInt(scoreRaw)
		if !ok {
			continue
		}
		result[teamID] = scoreValue
	}
	return result
}

func toStringAnyMap(raw any) (map[string]any, bool) {
	switch typed := raw.(type) {
	case map[string]any:
		return typed, true
	default:
		return nil, false
	}
}

func mapGetCaseInsensitive(m map[string]any, key string) (any, bool) {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return nil, false
}

func toInt(raw any) (int, bool) {
	switch n := raw.(type) {
	case float64:
		return int(n), true
	case float32:
		return int(n), true
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case json.Number:
		value, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int(value), true
	default:
		return 0, false
	}
}

func toString(raw any) string {
	switch v := raw.(type) {
	case string:
		return v
	default:
		return ""
	}
}
