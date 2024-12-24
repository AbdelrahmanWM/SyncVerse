package data_structure

import (
	event "github.com/AbdelrahmanWM/SyncVerse/crdt/Event"
	a "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
)

type CRDTDataStructure interface {
	Apply(e event.Event) // checks the event type, and applies the needed action (Insert/Delete/Update/etc)
	GetEvent(a a.Action) event.Event
}
