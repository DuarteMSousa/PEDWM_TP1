// Package services implements the application services layer of the system.
// Each service orchestrates business operations by delegating validation to the domain and persistence to the injected repositories.
//
// Available services:
//   - UserService: registration, authentication, and user queries.
//   - RoomService: room management (create, join, leave, start game).
//   - GameService: query and update game state.
//   - UserStatsService: user statistics (ELO, wins).
//   - EventService: persistence and query of domain events.
package services
