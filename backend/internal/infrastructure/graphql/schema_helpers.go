package graphqltransport

import (
	"backend/internal/application/ports"
	"fmt"

	gql "github.com/graphql-go/graphql"
)

func resolveRoomSource(source any) (ports.Room, error) {
	room, ok := source.(ports.Room)
	if ok {
		return room, nil
	}

	roomPtr, ok := source.(*ports.Room)
	if ok && roomPtr != nil {
		return *roomPtr, nil
	}

	return ports.Room{}, fmt.Errorf("unsupported room source type: %T", source)
}

func roomHostResolver(p gql.ResolveParams) (any, error) {
	room, err := resolveRoomSource(p.Source)
	if err != nil {
		return nil, err
	}
	return room.HostID, nil
}

func roomStatusResolver(p gql.ResolveParams) (any, error) {
	room, err := resolveRoomSource(p.Source)
	if err != nil {
		return nil, err
	}
	return string(room.Status), nil
}

func roomCreatedAtResolver(p gql.ResolveParams) (any, error) {
	room, err := resolveRoomSource(p.Source)
	if err != nil {
		return nil, err
	}
	return room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}

func toInt(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}
