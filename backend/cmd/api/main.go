package main

import (
	"backend/internal/application/ports"
	"backend/internal/application/usecases/players"
	"backend/internal/application/usecases/rooms"
	domainevents "backend/internal/domain/events"
	infraevents "backend/internal/infrastructure/events"
	graphqltransport "backend/internal/infrastructure/graphql"
	postgresstore "backend/internal/infrastructure/persistence/postgres"
	wstransport "backend/internal/infrastructure/websocket"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	hub := wstransport.NewHub()
	eventBus := domainevents.NewEventBus()
	_ = eventBus.Subscribe(wstransport.NewWebSocketObserver(hub))

	playerRepo, roomRepo, cleanup := buildPersistence()
	defer cleanup()

	playerUsecase := players.NewService(playerRepo)
	roomUsecase := rooms.NewService(roomRepo, infraevents.NewEventBusPublisher(eventBus))

	graphQLHandler, err := graphqltransport.NewHandler(playerUsecase, roomUsecase)
	if err != nil {
		log.Fatalf("failed to build graphql handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ws", wstransport.NewHandler(hub))
	mux.Handle("/graphql", graphQLHandler)

	addr := os.Getenv("API_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = ":4000"
	}

	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func buildPersistence() (ports.PlayerRepository, ports.RoomRepository, func()) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal(errors.New("DATABASE_URL is required (memory mode has been removed)"))
	}

	store, err := postgresstore.NewLobbyStore(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to initialize postgres store: %v", err)
	}

	log.Printf("persistence mode: postgres")
	return store, store, func() {
		store.Close()
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
