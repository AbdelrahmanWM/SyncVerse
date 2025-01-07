package event

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	vc "github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
)

type Event struct {
	Kind        EventType
	UserID      global.UserID
	ReplicaID   global.ReplicaID
	VectorClock vc.VectorClock
	Metadata    EventMetadata
}

func NewEvent(kind EventType, userID global.UserID, replicaID global.ReplicaID, vectorClock vc.VectorClock, metadata EventMetadata) *Event {
	return &Event{kind, userID, replicaID, vectorClock, metadata} // new event
}

// The function returns true of the called event have priority over the passed event, false otherwise
func (e *Event) Before(e2 *Event) bool {
	compare := e.VectorClock.Compare(e2.VectorClock)
	switch compare {
	case -1:
		return true
	case 1:
		return false
	case 0:
		hash := e.VectorClock.CompareHashes(e2.VectorClock)
		switch hash {
		case -1:
			return true
		case 0: /// should never happen
			fmt.Println("SAME EVENT!") // temp for debugging
			return false
		case 1:
			return false
		}
	}
	return false // shouldn't be reached
}
func (e *Event) String() string {
	var result strings.Builder
	result.WriteString("Event\n")
	result.WriteString(e.Kind.String())
	result.WriteString("\n")
	result.WriteString(string(e.UserID))
	result.WriteString("\n")
	result.WriteString(string(e.ReplicaID))
	result.WriteString("\n")
	result.WriteString(e.VectorClock.String())
	result.WriteString("\n")
	if e.Metadata != nil {
		result.WriteString(e.Metadata.String())
	} else {
		result.WriteString("No metadata provided")
	}
	result.WriteString("\n")
	return result.String()
}
func (e *Event) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(e)
}
func (e *Event) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(e)
}
