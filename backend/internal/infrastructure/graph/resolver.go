package graph

import application "backend/internal/application/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	RoomService       *application.RoomService
	UserService       *application.UserService
	FriendshipService *application.FriendshipService
	UserStatsService  *application.UserStatsService
}
