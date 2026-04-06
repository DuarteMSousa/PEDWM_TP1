package graph

import "backend/internal/application/services"

// Resolver groups the dependencies of the GraphQL resolvers.
// It serves as a dependency injection point for all resolvers.
type Resolver struct {
	RoomService       *services.RoomService
	UserService       *services.UserService
	FriendshipService *services.FriendshipService
	UserStatsService  *services.UserStatsService
	EventService      *services.EventService
	GameService       *services.GameService
}
