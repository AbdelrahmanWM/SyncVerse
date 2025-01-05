package event

import (
	"strings"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	vc "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

type Event struct {
	Kind        EventType
	UserID      global.UserID
	ReplicaID   global.ReplicaID
	VectorClock vc.VectorClock
	Metadata    EventMetadata
}

var EventMetadataRegistry map[EventType]EventMetadataConstructor = make(map[EventType]EventMetadataConstructor)

func NewEvent(kind EventType, userID global.UserID, replicaID global.ReplicaID, vectorClock vc.VectorClock, metadata EventMetadata) *Event {
	return &Event{kind, userID, replicaID, vectorClock, metadata} // new event
}

type EventMetadataConstructor func(inputs ...any) EventMetadata

type EventType int

const (
	Insert EventType = iota
	Delete
)
func (et *EventType) String()string{
	switch *et{
	case Insert:
		return "Insert"
	case Delete:
		return "Delete"
	default:
		return "Undefined"
	}
}
func registryNewEventMetadata(eventType EventType, eventMetadata EventMetadataConstructor) bool {
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
func (e *Event)String()string{
	var result strings.Builder
	result.WriteString("Event\n")
	result.WriteString(e.Kind.String())
	result.WriteString("\n")
	result.WriteString(string(e.UserID))
	result.WriteString("\n")
	result.WriteString(string(e.ReplicaID))
	result.WriteString("\n")
	result.WriteString(e.Metadata.String())
	result.WriteString("\n")
	return result.String()
}
func init() {
	initializeEventMetadataRegistry()
}
