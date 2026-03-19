package graphqltransport

import (
	"net/http"

	playersuc "backend/internal/application/usecases/players"
	roomsuc "backend/internal/application/usecases/rooms"

	gqlhandler "github.com/graphql-go/handler"
)

func NewHandler(playersService *playersuc.Service, roomsService *roomsuc.Service) (http.Handler, error) {
	schema, err := newSchema(playersService, roomsService)
	if err != nil {
		return nil, err
	}

	return gqlhandler.New(&gqlhandler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}), nil
}
