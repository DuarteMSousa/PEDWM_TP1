package graphqltransport

import (
	playersuc "backend/internal/application/usecases/players"
	roomsuc "backend/internal/application/usecases/rooms"

	gql "github.com/graphql-go/graphql"
)

func newSchema(playersService *playersuc.Service, roomsService *roomsuc.Service) (gql.Schema, error) {
	objects := newSchemaObjects(playersService)
	createRoomInput := newCreateRoomInput()
	queryType := newQueryType(objects, playersService, roomsService)
	mutationType := newMutationType(objects, createRoomInput, playersService, roomsService)

	return gql.NewSchema(gql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
