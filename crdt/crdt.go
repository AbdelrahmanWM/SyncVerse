package crdt

import (
	a "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	d "github.com/AbdelrahmanWM/SyncVerse/crdt/data_structure"
	e "github.com/AbdelrahmanWM/SyncVerse/crdt/event"
	v "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/utility/error_formatter.go"
)

type CRDT struct {
	replicaID          string
	currentVectorClock v.VectorClock
	dataStructure      d.CRDTDataStructure
	eventQueue         []*e.Event
}

func NewCRDT(dataStructure d.CRDTDataStructure, replicaID string) *CRDT {
	return &CRDT{
		replicaID,
		v.NewVectorClock(replicaID),
		dataStructure,
		make([]*e.Event, 10), //for now
	}
}

func (crdt *CRDT) Prepare(action *a.Action) (*e.Event, error) {
	event, err := crdt.dataStructure.GetEvent(action, crdt.currentVectorClock.NewVectorClock(crdt.replicaID))
	if err != nil || event == nil {
		return nil, err
	}
	crdt.eventQueue = append(crdt.eventQueue, event) // temp
	return event, nil                                ////temp
}
func (crdt *CRDT) Apply(event *e.Event) error {
	if !crdt.currentVectorClock.IsValidSuccessor(event.VectorClock) {
		return error_formatter.NewError("Missing event(s)") // temp
	}
	err := crdt.dataStructure.Apply(event)
	if err != nil {
		return err
	}
	crdt.currentVectorClock = crdt.currentVectorClock.Merge(event.VectorClock)
	return nil
}

func (crdt *CRDT) Query(remoteReplicaID string) []e.Event { // temp
	// webRTC interaction
	// todo: implementation
	return []e.Event{}
}
