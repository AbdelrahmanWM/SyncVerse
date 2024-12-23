package data_structure

import (
	a "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	e "github.com/AbdelrahmanWM/SyncVerse/crdt/event"
)

type CRDTDataStructure interface {
	Apply(e e.Event) // checks the event type, and applies the needed action (Insert/Delete/Update/etc)
	GetEvent(a a.Action) e.Event
}
