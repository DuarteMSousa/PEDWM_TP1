package graph

import "backend/internal/application/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	RoomService       *services.RoomService
	UserService       *services.UserService
	FriendshipService *services.FriendshipService
	UserStatsService  *services.UserStatsService
	EventService      *services.EventService
}
