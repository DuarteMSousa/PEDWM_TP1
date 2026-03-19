package graphqltransport

import (
	"errors"
	"fmt"

	lobbyapp "backend/internal/application/lobby"

	gql "github.com/graphql-go/graphql"
)

func newSchema(service *lobbyapp.Service) (gql.Schema, error) {
	playerType := gql.NewObject(gql.ObjectConfig{
		Name: "Player",
		Fields: gql.Fields{
			"id": &gql.Field{Type: gql.NewNonNull(gql.ID)},
			"nickname": &gql.Field{
				Type: gql.NewNonNull(gql.String),
			},
			"createdAt": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (any, error) {
					player, ok := p.Source.(lobbyapp.Player)
					if !ok {
						playerPtr, okPtr := p.Source.(*lobbyapp.Player)
						if !okPtr || playerPtr == nil {
							return "", nil
						}
						player = *playerPtr
					}
					return player.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
				},
			},
		},
	})

	roomStatusEnum := gql.NewEnum(gql.EnumConfig{
		Name: "RoomStatus",
		Values: gql.EnumValueConfigMap{
			"OPEN":    &gql.EnumValueConfig{Value: string(lobbyapp.RoomStatusOpen)},
			"IN_GAME": &gql.EnumValueConfig{Value: string(lobbyapp.RoomStatusInGame)},
			"CLOSED":  &gql.EnumValueConfig{Value: string(lobbyapp.RoomStatusClosed)},
		},
	})

	roomType := gql.NewObject(gql.ObjectConfig{
		Name: "Room",
		Fields: gql.Fields{
			"id":           &gql.Field{Type: gql.NewNonNull(gql.ID)},
			"name":         &gql.Field{Type: gql.NewNonNull(gql.String)},
			"hostPlayerId": &gql.Field{Type: gql.NewNonNull(gql.ID), Resolve: roomHostResolver},
			"status":       &gql.Field{Type: gql.NewNonNull(roomStatusEnum), Resolve: roomStatusResolver},
			"maxPlayers":   &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"isPrivate":    &gql.Field{Type: gql.NewNonNull(gql.Boolean)},
			"createdAt": &gql.Field{
				Type:    gql.NewNonNull(gql.String),
				Resolve: roomCreatedAtResolver,
			},
			"playerIds": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.ID))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					room, err := resolveRoomSource(p.Source)
					if err != nil {
						return nil, err
					}
					return room.PlayerIDs, nil
				},
			},
			"playersCount": &gql.Field{
				Type: gql.NewNonNull(gql.Int),
				Resolve: func(p gql.ResolveParams) (any, error) {
					room, err := resolveRoomSource(p.Source)
					if err != nil {
						return nil, err
					}
					return len(room.PlayerIDs), nil
				},
			},
			"phase": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (any, error) {
					room, err := resolveRoomSource(p.Source)
					if err != nil {
						return nil, err
					}
					switch room.Status {
					case lobbyapp.RoomStatusInGame:
						return "PLAYING_TRICK", nil
					case lobbyapp.RoomStatusClosed:
						return "ENDED", nil
					default:
						return "WAITING_FOR_PLAYERS", nil
					}
				},
			},
			"trumpSuit": &gql.Field{
				Type: gql.String,
				Resolve: func(p gql.ResolveParams) (any, error) {
					_, err := resolveRoomSource(p.Source)
					if err != nil {
						return nil, err
					}
					return "HEARTS", nil
				},
			},
			"players": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(playerType))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					room, err := resolveRoomSource(p.Source)
					if err != nil {
						return nil, err
					}
					return service.PlayersByIDs(room.PlayerIDs), nil
				},
			},
		},
	})

	profileType := gql.NewObject(gql.ObjectConfig{
		Name: "Profile",
		Fields: gql.Fields{
			"id":       &gql.Field{Type: gql.NewNonNull(gql.ID)},
			"nickname": &gql.Field{Type: gql.NewNonNull(gql.String)},
			"matches":  &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"wins":     &gql.Field{Type: gql.NewNonNull(gql.Int)},
		},
	})

	createRoomInput := gql.NewInputObject(gql.InputObjectConfig{
		Name: "CreateRoomInput",
		Fields: gql.InputObjectConfigFieldMap{
			"name": &gql.InputObjectFieldConfig{
				Type: gql.NewNonNull(gql.String),
			},
			"hostPlayerId": &gql.InputObjectFieldConfig{
				Type: gql.NewNonNull(gql.ID),
			},
			"maxPlayers": &gql.InputObjectFieldConfig{
				Type: gql.Int,
			},
			"isPrivate": &gql.InputObjectFieldConfig{
				Type: gql.Boolean,
			},
		},
	})

	queryType := gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"players": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(playerType))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					return service.ListPlayers(), nil
				},
			},
			"rooms": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(roomType))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					return service.ListRoomsDetailed(), nil
				},
			},
			"room": &gql.Field{
				Type: roomType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					id, _ := p.Args["id"].(string)
					room, ok := service.GetRoom(id)
					if !ok {
						return nil, nil
					}
					return room, nil
				},
			},
			"profile": &gql.Field{
				Type: gql.NewNonNull(profileType),
				Args: gql.FieldConfigArgument{
					"userId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					userID, _ := p.Args["userId"].(string)
					player, ok := service.GetPlayer(userID)
					if ok {
						return map[string]any{
							"id":       player.ID,
							"nickname": player.Nickname,
							"matches":  24,
							"wins":     14,
						}, nil
					}

					return map[string]any{
						"id":       userID,
						"nickname": userID,
						"matches":  0,
						"wins":     0,
					}, nil
				},
			},
		},
	})

	mutationType := gql.NewObject(gql.ObjectConfig{
		Name: "Mutation",
		Fields: gql.Fields{
			"createPlayer": &gql.Field{
				Type: gql.NewNonNull(playerType),
				Args: gql.FieldConfigArgument{
					"nickname": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					nickname, _ := p.Args["nickname"].(string)
					return service.CreatePlayer(nickname)
				},
			},
			"createRoom": &gql.Field{
				Type: gql.NewNonNull(roomType),
				Args: gql.FieldConfigArgument{
					"input": &gql.ArgumentConfig{Type: gql.NewNonNull(createRoomInput)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					rawInput, ok := p.Args["input"].(map[string]any)
					if !ok {
						return nil, errors.New("input is required")
					}

					name, _ := rawInput["name"].(string)
					hostPlayerID, _ := rawInput["hostPlayerId"].(string)
					maxPlayers := toInt(rawInput["maxPlayers"])
					isPrivate, _ := rawInput["isPrivate"].(bool)

					return service.CreateRoom(name, hostPlayerID, maxPlayers, isPrivate)
				},
			},
			"deleteRoom": &gql.Field{
				Type: gql.NewNonNull(gql.Boolean),
				Args: gql.FieldConfigArgument{
					"roomId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"requesterId": &gql.ArgumentConfig{
						Type: gql.NewNonNull(gql.ID),
					},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					requesterID, _ := p.Args["requesterId"].(string)
					if err := service.DeleteRoom(roomID, requesterID); err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"joinRoom": &gql.Field{
				Type: gql.NewNonNull(roomType),
				Args: gql.FieldConfigArgument{
					"roomId":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"playerId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					playerID, _ := p.Args["playerId"].(string)
					if _, err := service.JoinRoom(roomID, playerID); err != nil {
						return nil, err
					}
					room, _ := service.GetRoom(roomID)
					return room, nil
				},
			},
			"leaveRoom": &gql.Field{
				Type: gql.NewNonNull(roomType),
				Args: gql.FieldConfigArgument{
					"roomId":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"playerId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					playerID, _ := p.Args["playerId"].(string)
					if _, err := service.LeaveRoom(roomID, playerID); err != nil {
						return nil, err
					}
					room, _ := service.GetRoom(roomID)
					return room, nil
				},
			},
		},
	})

	return gql.NewSchema(gql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

func resolveRoomSource(source any) (lobbyapp.Room, error) {
	room, ok := source.(lobbyapp.Room)
	if ok {
		return room, nil
	}

	roomPtr, ok := source.(*lobbyapp.Room)
	if ok && roomPtr != nil {
		return *roomPtr, nil
	}

	return lobbyapp.Room{}, fmt.Errorf("unsupported room source type: %T", source)
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
