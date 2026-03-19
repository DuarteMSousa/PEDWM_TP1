package main

import (
	"backend/internal/application/usecases/players"
	"backend/internal/application/usecases/rooms"
	domainevents "backend/internal/domain/events"
	infraevents "backend/internal/infrastructure/events"
	graphqltransport "backend/internal/infrastructure/graphql"
	"backend/internal/infrastructure/persistence/memory"
	wstransport "backend/internal/infrastructure/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	hub := wstransport.NewHub()
	eventBus := domainevents.NewEventBus()
	_ = eventBus.Subscribe(wstransport.NewWebSocketObserver(hub))

	store := memory.NewLobbyStore()
	playerUsecase := players.NewService(store)
	roomUsecase := rooms.NewService(store, infraevents.NewEventBusPublisher(eventBus))

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
