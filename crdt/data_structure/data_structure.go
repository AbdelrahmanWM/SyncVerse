package data_structure

import (
	a "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	event "github.com/AbdelrahmanWM/SyncVerse/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type CRDTDataStructure interface {
	Apply(e *event.Event)error // checks the event type, and applies the needed action (Insert/Delete/Update/etc)
	GetEvent(a *a.Action, vectorClock vector_clock.VectorClock) (*event.Event,error)
}
