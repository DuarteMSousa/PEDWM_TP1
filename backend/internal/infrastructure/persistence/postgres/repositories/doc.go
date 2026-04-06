// Package repositories contains the concrete implementations of the repositories
// defined in application/interfaces, using PostgreSQL (pgx/v5) as the
// persistence mechanism. Each repository operates on a shared connection pool
// and uses transactions when necessary to ensure atomicity.
//
// Available repositories:
//   - UserPostgresRepository: user persistence.
//   - UserStatsPostgresRepository: statistics persistence.
//   - RoomPostgresRepository: room and room players persistence.
//   - GamePostgresRepository: game persistence.
//   - FriendshipPostgresRepository: friendship persistence.
//   - EventPostgresRepository: domain event persistence.
package repositories
