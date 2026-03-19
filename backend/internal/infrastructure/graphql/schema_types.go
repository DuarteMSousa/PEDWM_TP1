package graphqltransport

import (
	"backend/internal/application/ports"
	playersuc "backend/internal/application/usecases/players"

	gql "github.com/graphql-go/graphql"
)

type schemaObjects struct {
	playerType     *gql.Object
	roomStatusEnum *gql.Enum
	roomType       *gql.Object
	profileType    *gql.Object
}

func newSchemaObjects(playersService *playersuc.Service) schemaObjects {
	playerType := gql.NewObject(gql.ObjectConfig{
		Name: "Player",
		Fields: gql.Fields{
			"id":       &gql.Field{Type: gql.NewNonNull(gql.ID)},
			"nickname": &gql.Field{Type: gql.NewNonNull(gql.String)},
			"createdAt": &gql.Field{
				Type: gql.NewNonNull(gql.String),
				Resolve: func(p gql.ResolveParams) (any, error) {
					player, ok := p.Source.(ports.Player)
					if !ok {
						playerPtr, okPtr := p.Source.(*ports.Player)
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
			"OPEN":    &gql.EnumValueConfig{Value: string(ports.RoomStatusOpen)},
			"IN_GAME": &gql.EnumValueConfig{Value: string(ports.RoomStatusInGame)},
			"CLOSED":  &gql.EnumValueConfig{Value: string(ports.RoomStatusClosed)},
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
					case ports.RoomStatusInGame:
						return "PLAYING_TRICK", nil
					case ports.RoomStatusClosed:
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
					return playersService.PlayersByIDs(room.PlayerIDs), nil
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

	return schemaObjects{
		playerType:     playerType,
		roomStatusEnum: roomStatusEnum,
		roomType:       roomType,
		profileType:    profileType,
	}
}

func newCreateRoomInput() *gql.InputObject {
	return gql.NewInputObject(gql.InputObjectConfig{
		Name: "CreateRoomInput",
		Fields: gql.InputObjectConfigFieldMap{
			"name":         &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"hostPlayerId": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.ID)},
			"maxPlayers":   &gql.InputObjectFieldConfig{Type: gql.Int},
			"isPrivate":    &gql.InputObjectFieldConfig{Type: gql.Boolean},
			"password":     &gql.InputObjectFieldConfig{Type: gql.String},
		},
	})
}
