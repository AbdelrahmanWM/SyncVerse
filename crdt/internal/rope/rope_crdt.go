package rope

import (
	action "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	event "github.com/AbdelrahmanWM/SyncVerse/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/value"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/utility/error_formatter.go"
)

func (r *Rope) GetEvent(a *action.Action, vectorClock vector_clock.VectorClock) (*event.Event, error) {
	var e event.Event
	e.UserID = a.UserID
	e.ReplicaID = a.ReplicaID
	e.VectorClock = vectorClock
	switch a.Kind {
	case action.Insert:
		actionMetadata, ok := a.Metadata.(action.InsertionMetadata)
		if !ok {
			return nil, error_formatter.NewError("Invalid insertion metadata")
		}
		e.Kind = event.Insert
		insertionClockOffset := r.findInsertionBlockOffset(actionMetadata.Index)
		if insertionClockOffset == nil {
			return nil, error_formatter.NewError("Insertion block now found")
		}
		contentBlock := block.NewBlock(vector_clock.NewClockOffset(vectorClock, 0), actionMetadata.Content, value.ByteBuffer, false)
		e.Metadata = event.NewInsertionEventMetadata(contentBlock, insertionClockOffset, actionMetadata.Index)
	case action.Delete:
		actionMetadata, ok := a.Metadata.(action.DeletionMetadata)
		if !ok {
			return nil, error_formatter.NewError("Invalid deletion metadata") // todo:Make it an error
		}
		e.Kind = event.Delete
		deletionMetadata := r.findBlocks(actionMetadata.Index, actionMetadata.Length)
		if deletionMetadata == nil {
			return nil, error_formatter.NewError("no deletion metadata found")
		}
		e.Metadata = event.NewDeletionEventMetadata(deletionMetadata, actionMetadata.Index)
	}
	return &e, nil
}

func (r *Rope) Apply(e *event.Event)error {
	switch e.Kind {
	case event.Insert:
		insertionMetadata, ok := e.Metadata.(event.InsertionEventMetadata)
		if !ok {
			return error_formatter.NewError("Invalid insertion metadata") //temp
		}
		inserted:=r.Insert(insertionMetadata.ContentBlock,insertionMetadata.ToBeInsertedAfter,insertionMetadata.StartIndex)
		if !inserted { /// temp
			return error_formatter.NewError("Failed to apply insertion event")
		}
	case event.Delete:
		deletionMetadata,ok:=e.Metadata.(event.DeletionEventMetadata)
		if !ok{
			return error_formatter.NewError("Invalid deletion metadata")
		}
		deleted:=r.Delete(deletionMetadata.DeletionMetadata,deletionMetadata.StartIndex)
		if !deleted {
			return error_formatter.NewError("Failed to apply deletion event")
		}
	}
	return nil
}
