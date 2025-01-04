package event

import (
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	vc "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

type Event struct {
	Kind        EventType
	UserID      global.UserID
	ReplicaID   global.ReplicaID
	VectorClock vc.VectorClock
	Metadata    any
}

var EventMetadataRegistry map[EventType]EventMetadata = make(map[EventType]EventMetadata)

func NewEvent(kind EventType, userID global.UserID, replicaID global.ReplicaID, vectorClock vc.VectorClock, metadata any) *Event {
	return &Event{kind, userID, replicaID, vectorClock, metadata} // new event
}

type EventMetadata func(inputs ...any) any

type EventType int

const (
	Insert EventType = iota
	Delete
)

func registryNewEventMetadata(eventType EventType, eventMetadata EventMetadata) bool {
	_, ok := EventMetadataRegistry[eventType]
	if ok || eventMetadata == nil {
		return false
	}
	EventMetadataRegistry[eventType] = eventMetadata
	return true
}
func initializeEventMetadataRegistry() {
	registryNewEventMetadata(Insert, NewInsertionEventMetadata)
	registryNewEventMetadata(Delete, NewDeletionEventMetadata)
}
func init() {
	initializeEventMetadataRegistry()
}
