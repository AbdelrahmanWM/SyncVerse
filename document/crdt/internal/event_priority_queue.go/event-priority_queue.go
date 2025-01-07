package event_priority_queue

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/document/error_formatter.go"
)

type EventPriorityQueue struct {
	eventHeap *EventHeap
}

func NewEventPriorityQueue(items []*event.Event) *EventPriorityQueue {
	var mux sync.RWMutex
	eh := EventHeap{[]*EventItem{}, &mux}
	eh.items = make([]*EventItem, len(items))
	for i, item := range items {
		eh.items[i] = &EventItem{item, i}
	}

	heap.Init(&eh)
	return &EventPriorityQueue{&eh}
}

func (epq *EventPriorityQueue) Push(e *event.Event) {
	eventItem := &EventItem{
		e,
		0, // will be updated on Push
	}
	heap.Push(epq.eventHeap, eventItem)
}
func (epq *EventPriorityQueue) Pop() (*event.Event, error) {
	item, ok := heap.Pop(epq.eventHeap).(*EventItem)
	if !ok {
		return nil, error_formatter.NewError("invalid item type")
	}
	event := item.value
	return event, nil
}
func (epq *EventPriorityQueue) Len() int {
	return epq.eventHeap.Len()
}
func main() {
	items := []*event.Event{ // random events
		event.NewEvent(event.Insert, "C", "C", vector_clock.VectorClock{"A": 1, "C": 1}, event.NewInsertionEventMetadata()),
		event.NewEvent(event.Insert, "B", "B", vector_clock.VectorClock{"B": 1}, event.NewInsertionEventMetadata()),
		event.NewEvent(event.Insert, "A", "A", vector_clock.VectorClock{"A": 2, "B": 1}, event.NewInsertionEventMetadata()),
		event.NewEvent(event.Insert, "B", "B", vector_clock.VectorClock{}, event.NewInsertionEventMetadata()),
		event.NewEvent(event.Insert, "A", "A", vector_clock.VectorClock{"A": 1}, event.NewInsertionEventMetadata()),
	}
	epq := NewEventPriorityQueue(items)

	e := event.NewEvent(event.Insert, "A", "A", vector_clock.VectorClock{"A": 3, "B": 2, "C": 1}, event.NewInsertionEventMetadata())
	epq.Push(e)
	e = event.NewEvent(event.Insert, "A", "A", vector_clock.VectorClock{"A": 1, "B": 2, "C": 2}, event.NewInsertionEventMetadata())
	epq.Push(e)
	for epq.Len() > 0 {
		item, err := epq.Pop()
		if err != nil {
			fmt.Println("ERROR")
			break
		}
		fmt.Println(item.String())
	}
}
