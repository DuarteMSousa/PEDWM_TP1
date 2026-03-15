package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"backend/internal/graph"
	"backend/internal/graph/generated"
	"backend/internal/repository"
	"backend/internal/service"
)

func main() {

	roomRepo := repository.NewRoomRepository()
	roomService := service.NewRoomService(roomRepo)

	resolver := &graph.Resolver{
		RoomService: roomService,
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Println("server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
