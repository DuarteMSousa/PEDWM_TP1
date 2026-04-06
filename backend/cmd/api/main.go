// The main package is the entry point for the Sueca game's backend application.
// It configures and initializes all system components: database, event bus,
// WebSocket hub, application services, GraphQL, and HTTP server.
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
	"log/slog"
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
	// Configure the global structured logger.
	logLevel := slog.LevelInfo

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	slog.SetDefault(logger)

	ctx := context.Background()

	// ========================
	// Base infrastructure
	// ========================
	slog.Info("initializing base infrastructure")

	hub := wstransport.GetHubInstance()

	eventBus := events.NewEventBus()
	events.SetDefaultBus(eventBus)
	eventBus.Subscribe(wstransport.NewWebSocketObserver(hub))

	// ========================
	// Persistence
	// ========================
	slog.Info("establishing database connection")

	pool, err := postgres.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := postgres.EnsureSchema(ctx, pool); err != nil {
		slog.Error("failed to ensure database schema", "error", err)
		os.Exit(1)
	}
	slog.Info("database schema verified successfully")

	repo := repositories.NewRoomPostgresRepository(pool)
	userRepo := repositories.NewUserPostgresRepository(pool)
	friendshipRepo := repositories.NewFriendshipPostgresRepository(pool)
	userStatsRepo := repositories.NewUserStatsPostgresRepository(pool)
	gameRepo := repositories.NewGamePostgresRepository(pool)
	eventRepo := repositories.NewEventPostgresRepository(pool)

	// ========================
	// Command Dispatcher
	// ========================
	slog.Info("registering command handlers")

	dispatcher := wstransport.GetCommandDispatcherInstance()
	dispatcher.Register("play_card", wstransport.NewPlayCardHandler(hub))
	dispatcher.Register("change_bot_strategy", wstransport.NewChangeBotStrategyHandler(hub))

	// ========================
	// Application
	// ========================
	slog.Info("initializing application services")

	eventService := services.NewEventService(eventRepo)
	roomService := services.NewRoomService(repo, gameRepo, userRepo, eventService, hub)
	userService := services.NewUserService(userRepo, userStatsRepo)
	friendshipService := services.NewFriendshipService(friendshipRepo, userRepo)
	userStatsService := services.NewUserStatsService(userStatsRepo, userRepo)
	gameService := services.NewGameService(gameRepo)

	// ========================
	// Event Dispatcher
	// ========================
	slog.Info("registering event handlers")

	eventDispatcher := events_infrastructure.GetEventDispatcherInstance()
	// eventDispatcher.Register(string(events.EventPlayerLeft), events_infrastructure.NewPlayerLeftEventHandler(roomService))
	eventDispatcher.Register(string(events.EventGameEnded), events_infrastructure.NewGameEndedEventHandler(userStatsService, gameService))
	eventDispatcher.Register(string(events.EventRoomClosed), events_infrastructure.NewRoomClosedEventHandler(roomService))

	// ========================
	// GraphQL
	// ========================
	slog.Info("configuring GraphQL server")

	resolver := &graph.Resolver{
		RoomService:       roomService,
		UserService:       userService,
		FriendshipService: friendshipService,
		UserStatsService:  userStatsService,
		EventService:      eventService,
		GameService:       gameService,
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

	mux.Handle("/ws", wstransport.NewHandler(hub, dispatcher, roomService))
	mux.Handle("/graphql", srv)
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	addr := os.Getenv("API_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = ":7000"
	}

	slog.Info("HTTP server starting", "addr", addr)

	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		slog.Error("HTTP server terminated with error", "error", err)
		os.Exit(1)
	}
}

// withCORS wraps an http.Handler with permissive CORS headers.
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
