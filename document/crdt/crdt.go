package crdt

import (
	action "github.com/AbdelrahmanWM/SyncVerse/document/crdt/action"
	"github.com/AbdelrahmanWM/SyncVerse/document/error_formatter.go"

	d "github.com/AbdelrahmanWM/SyncVerse/document/crdt/data_structure"
	event "github.com/AbdelrahmanWM/SyncVerse/document/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/value"
	v "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

type CRDT struct {
	replicaID          string
	currentVectorClock v.VectorClock
	dataStructure      d.CRDTDataStructure
	eventQueue         []*event.Event
}

func NewCRDT(dataStructure d.CRDTDataStructure, replicaID string) *CRDT {
	return &CRDT{
		replicaID,
		v.NewVectorClock(replicaID),
		dataStructure,
		make([]*event.Event, 10), //for now
	}
}

func (crdt *CRDT) Prepare(a *action.Action) (*event.Event, error) {
	newVectorClock := crdt.currentVectorClock.NewVectorClock(crdt.replicaID)
	var e event.Event
	e.UserID = a.UserID
	e.ReplicaID = a.ReplicaID
	e.VectorClock = newVectorClock
	switch a.Kind {
	case action.Insert:
		actionMetadata, ok := a.Metadata.(action.InsertionMetadata)
		if !ok {
			return nil, error_formatter.NewError("Invalid insertion metadata")
		}
		e.Kind = event.Insert
		insertionClockOffset, err := crdt.dataStructure.FindInsertionBlockOffset(actionMetadata.Index)
		if err != nil {
			return nil, err
		}
		contentBlock := block.NewBlock(v.NewClockOffset(newVectorClock, 0), actionMetadata.Content, value.ByteBuffer, false)
		e.Metadata = event.NewInsertionEventMetadata(contentBlock, insertionClockOffset, actionMetadata.Index)
	case action.Delete:
		actionMetadata, ok := a.Metadata.(action.DeletionMetadata)
		if !ok {
			return nil, error_formatter.NewError("Invalid deletion metadata") // todo:Make it an error
		}
		e.Kind = event.Delete
		deletionMetadata, err := crdt.dataStructure.FindBlocks(actionMetadata.Index, actionMetadata.Length)
		if err != nil {
			return nil, err
		}
		e.Metadata = event.NewDeletionEventMetadata(deletionMetadata, actionMetadata.Index)
	default:
		return nil, error_formatter.NewError("Action type not found")
	}
	crdt.eventQueue = append(crdt.eventQueue, &e) // temp
	return &e, nil
}
func (crdt *CRDT) Apply(e *event.Event) error {
	if !crdt.currentVectorClock.IsValidSuccessor(e.VectorClock) {
		return error_formatter.NewError("Missing event(s)") // temp
	}
	switch e.Kind {
	case event.Insert:
		insertionMetadata, ok := e.Metadata.(event.InsertionEventMetadata)
		if !ok {
			return error_formatter.NewError("Invalid insertion metadata")
		}
		err := crdt.dataStructure.Insert(insertionMetadata.ContentBlock, insertionMetadata.ToBeInsertedAfter, insertionMetadata.StartIndex)
		if err != nil {
			return error_formatter.NewError("Failed to apply insertion event")
		}
	case event.Delete:
		deletionMetadata, ok := e.Metadata.(event.DeletionEventMetadata)
		if !ok {
			return error_formatter.NewError("Invalid deletion metadata")
		}
		err := crdt.dataStructure.Delete(deletionMetadata.DeletionMetadata, deletionMetadata.StartIndex)
		if err != nil {
			return error_formatter.NewError("Failed to apply deletion event")
		}
	default:
		return error_formatter.NewError("Event type not found")
	}
	crdt.currentVectorClock = crdt.currentVectorClock.Merge(e.VectorClock) // merging the latest
	return nil
}

func (crdt *CRDT) Query(remoteReplicaID string) []event.Event { // temp
	// webRTC interaction
	// todo: implementation
	return []event.Event{}
}
