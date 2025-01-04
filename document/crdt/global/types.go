package global

import "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"

type ReplicaID string
type UserID string
type ModifyMetadata struct {
	ClockOffset *vector_clock.ClockOffset
	Rng         [2]int
}
