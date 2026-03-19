package graphqltransport

import (
	"errors"

	playersuc "backend/internal/application/usecases/players"
	roomsuc "backend/internal/application/usecases/rooms"

	gql "github.com/graphql-go/graphql"
)

func newMutationType(
	objects schemaObjects,
	createRoomInput *gql.InputObject,
	playersService *playersuc.Service,
	roomsService *roomsuc.Service,
) *gql.Object {
	return gql.NewObject(gql.ObjectConfig{
		Name: "Mutation",
		Fields: gql.Fields{
			"createPlayer": &gql.Field{
				Type: gql.NewNonNull(objects.playerType),
				Args: gql.FieldConfigArgument{
					"nickname": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					nickname, _ := p.Args["nickname"].(string)
					return playersService.CreatePlayer(nickname)
				},
			},
			"createRoom": &gql.Field{
				Type: gql.NewNonNull(objects.roomType),
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
					password, _ := rawInput["password"].(string)

					return roomsService.CreateRoom(name, hostPlayerID, maxPlayers, isPrivate, password)
				},
			},
			"deleteRoom": &gql.Field{
				Type: gql.NewNonNull(gql.Boolean),
				Args: gql.FieldConfigArgument{
					"roomId":      &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"requesterId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					requesterID, _ := p.Args["requesterId"].(string)
					if err := roomsService.DeleteRoom(roomID, requesterID); err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"joinRoom": &gql.Field{
				Type: gql.NewNonNull(objects.roomType),
				Args: gql.FieldConfigArgument{
					"roomId":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"playerId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"password": &gql.ArgumentConfig{Type: gql.String},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					playerID, _ := p.Args["playerId"].(string)
					password, _ := p.Args["password"].(string)
					if _, err := roomsService.JoinRoom(roomID, playerID, password); err != nil {
						return nil, err
					}
					room, _ := roomsService.GetRoom(roomID)
					return room, nil
				},
			},
			"leaveRoom": &gql.Field{
				Type: gql.NewNonNull(objects.roomType),
				Args: gql.FieldConfigArgument{
					"roomId":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
					"playerId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					roomID, _ := p.Args["roomId"].(string)
					playerID, _ := p.Args["playerId"].(string)
					if _, err := roomsService.LeaveRoom(roomID, playerID); err != nil {
						return nil, err
					}
					room, _ := roomsService.GetRoom(roomID)
					return room, nil
				},
			},
		},
	})
}
