package format

import (
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/action"
)

type Format struct {
	Kind     ActionCode
	Metadata string // can be changed to any/map
}
