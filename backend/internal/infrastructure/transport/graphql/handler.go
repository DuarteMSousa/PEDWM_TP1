package graphqltransport

import (
	"net/http"

	lobbyapp "backend/internal/application/lobby"

	gqlhandler "github.com/graphql-go/handler"
)

func NewHandler(service *lobbyapp.Service) (http.Handler, error) {
	schema, err := newSchema(service)
	if err != nil {
		return nil, err
	}

	return gqlhandler.New(&gqlhandler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}), nil
}
