package graphqltransport

import (
	playersuc "backend/internal/application/usecases/players"
	roomsuc "backend/internal/application/usecases/rooms"

	gql "github.com/graphql-go/graphql"
)

func newQueryType(
	objects schemaObjects,
	playersService *playersuc.Service,
	roomsService *roomsuc.Service,
) *gql.Object {
	return gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"players": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(objects.playerType))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					return playersService.ListPlayers(), nil
				},
			},
			"rooms": &gql.Field{
				Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(objects.roomType))),
				Resolve: func(p gql.ResolveParams) (any, error) {
					return roomsService.ListRoomsDetailed(), nil
				},
			},
			"room": &gql.Field{
				Type: objects.roomType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					id, _ := p.Args["id"].(string)
					room, ok := roomsService.GetRoom(id)
					if !ok {
						return nil, nil
					}
					return room, nil
				},
			},
			"profile": &gql.Field{
				Type: gql.NewNonNull(objects.profileType),
				Args: gql.FieldConfigArgument{
					"userId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.ID)},
				},
				Resolve: func(p gql.ResolveParams) (any, error) {
					userID, _ := p.Args["userId"].(string)
					player, ok := playersService.GetPlayer(userID)
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
}
