package event

import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type Event struct {
	kind        EventType
	clockOffset *ClockOffset
	offset      int
	content     string
}

func NewEvent(kind EventType, vectorClock VectorClock, offset int, content string) *Event {
	return &Event{kind, NewClockOffset(vectorClock, offset), offset, content}
}

type EventType int

const (
	Insert EventType = iota
	Delete
)
