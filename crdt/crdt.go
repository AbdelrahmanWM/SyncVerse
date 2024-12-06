package crdt

import (
	a "github.com/AbdelrahmanWM/SyncVerse/crdt/Action"
	e "github.com/AbdelrahmanWM/SyncVerse/crdt/Event"
	d "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/data_structure"
	v "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type CRDT struct {
	replicaID          string
	currentVectorClock v.VectorClock
	dataStructure      d.CRDTDataStructure
	eventQueue         []e.Event
}

func NewCRDT(dataStructure d.CRDTDataStructure, replicaID string) *CRDT {
	return &CRDT{
		replicaID,
		v.NewVectorClock(replicaID),
		dataStructure,
		make([]e.Event, 10), //for now
	}
}

func (crdt *CRDT) Prepare(action a.Action) e.Event {
	event := crdt.dataStructure.GetEvent(action)
	crdt.eventQueue = append(crdt.eventQueue, event) // for now
	return event
}
func (crdt *CRDT) Effect(event e.Event) {
	crdt.dataStructure.Apply(event)
}

func (crdt *CRDT) Query(remoteReplicaID string)[]e.Event { // temp
	// webRTC interaction
	// todo: implementation 
	return []e.Event{} 
}
