package main

import (
	"backend/internal/application/services"
	"backend/internal/domain/events"
	events_infrastructure "backend/internal/infrastructure/events"
	"backend/internal/infrastructure/graph"
	"backend/internal/infrastructure/persistence/postgres"
	"backend/internal/infrastructure/persistence/postgres/repositories"
	wstransport "backend/internal/infrastructure/websocket"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	ctx := context.Background()

	// ========================
	// Infra base
	// ========================
	hub := wstransport.GetHubInstance()

	eventBus := events.NewEventBus()
	events.SetDefaultBus(eventBus)
	eventBus.Subscribe(wstransport.NewWebSocketObserver(hub))

	// ========================
	// Persistence
	// ========================
	pool, err := postgres.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := postgres.EnsureSchema(ctx, pool); err != nil {
		log.Fatal(err)
	}

	repo := repositories.NewRoomPostgresRepository(pool)
	userRepo := repositories.NewUserPostgresRepository(pool)
	friendshipRepo := repositories.NewFriendshipPostgresRepository(pool)
	userStatsRepo := repositories.NewUserStatsPostgresRepository(pool)
	gameRepo := repositories.NewGamePostgresRepository(pool)
	eventRepo := repositories.NewEventPostgresRepository(pool)

	// ========================
	// Command Dispatcher
	// ========================
	dispatcher := wstransport.GetCommandDispatcherInstance()
	dispatcher.Register("play_card", wstransport.NewPlayCardHandler(hub))

	// ========================
	// Application
	// ========================
	eventService := services.NewEventService(eventRepo)
	roomService := services.NewRoomService(repo, gameRepo, userRepo, eventService, hub)
	userService := services.NewUserService(userRepo, userStatsRepo)
	friendshipService := services.NewFriendshipService(friendshipRepo, userRepo)
	userStatsService := services.NewUserStatsService(userStatsRepo, userRepo)

	// ========================
	// Event Dispatcher
	// ========================
	eventDispatcher := events_infrastructure.GetEventDispatcherInstance()
	eventDispatcher.Register("PLAYER_LEFT", events_infrastructure.NewPlayerLeftEventHandler(roomService))
	eventDispatcher.Register("GAME_ENDED", events_infrastructure.NewGameEndedEventHandler(userStatsService))

	// ========================
	// GraphQL
	// ========================
	resolver := &graph.Resolver{
		RoomService:       roomService,
		UserService:       userService,
		FriendshipService: friendshipService,
		UserStatsService:  userStatsService,
		EventService:      eventService,
	}

	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: resolver},
	))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// ========================
	// HTTP
	// ========================
	mux := http.NewServeMux()

	mux.Handle("/ws", wstransport.NewHandler(hub, dispatcher))
	mux.Handle("/graphql", srv)
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	addr := os.Getenv("API_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = ":7000"
	}

	log.Printf("server running at http://localhost%s", addr)

	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
