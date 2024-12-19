package format
import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/action"
)
type Format struct {
	Kind ActionCode
	Metadata string // can be changed to any/map
}