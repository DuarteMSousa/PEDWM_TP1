// Package turnorder implements a deterministic circular queue to
// control the turn order of players. It ensures consistent ordering
// regardless of map iteration in Go, using Sequence and ID
// as tiebreaker criteria.
package turnorder
