package event_priority_queue

import (
	"container/heap"
	"sync"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/event"
)

type EventItem struct {
	value *event.Event
	index int
}

type EventHeap struct {
	items []*EventItem
	mux   *sync.RWMutex
}

func (eh *EventHeap) Len() int {
	eh.mux.RLock()
	defer eh.mux.RUnlock()
	return len(eh.items)
}

func (eh *EventHeap) Push(element any) {
	eh.mux.Lock()
	defer eh.mux.Unlock()
	n := len(eh.items)
	item := element.(*EventItem)
	item.index = n
	eh.items = append(eh.items, item)
}
func (eh *EventHeap) Pop() any {
	eh.mux.Lock()
	defer eh.mux.Unlock()
	oeh := eh.items
	n := len(oeh)
	item := oeh[n-1]
	oeh[n-1] = nil
	item.index = -1
	eh.items = oeh[0 : n-1]
	return item
}
func (eh *EventHeap) Swap(i, j int) {
	eh.mux.Lock()
	defer eh.mux.Unlock()
	eh.items[i], eh.items[j] = eh.items[j], eh.items[i]
	eh.items[i].index = i
	eh.items[j].index = j
}
func (eh *EventHeap) Less(i, j int) bool {
	eh.mux.RLock()
	defer eh.mux.RUnlock()
	return eh.items[i].value.Before(eh.items[j].value)
}
func (eh *EventHeap) update(eventItem *EventItem, value *event.Event) {
	eh.mux.Lock()
	defer eh.mux.Unlock()
	eventItem.value = value
	heap.Fix(eh, eventItem.index)
}