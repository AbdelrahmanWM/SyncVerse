package event
import (
	vc "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type Event struct {
	kind        EventType
	clockOffset *vc.ClockOffset
	offset      int
	content     string
}

func NewEvent(kind EventType, vectorClock vc.VectorClock, offset int, content string) *Event {
	return &Event{kind, vc.NewClockOffset(vectorClock, offset), offset, content}
}

type EventType int

const (
	Insert EventType = iota
	Delete
)
