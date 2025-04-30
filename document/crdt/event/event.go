package event

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	s "github.com/AbdelrahmanWM/SyncVerse/document/crdt/sequence_number"
	vc "github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/document/global"
)

type Event struct {
	Kind                 EventType
	OriginID             global.ReplicaID
	VectorClock          vc.VectorClock
	LocalSequenceNumber  s.SeqNum
	OriginSequenceNumber s.SeqNum
	Metadata             EventMetadata
}

func NewEvent(kind EventType, OriginID global.ReplicaID, vectorClock vc.VectorClock, localSequenceNumber s.SeqNum, originSequenceNumber s.SeqNum, metadata EventMetadata) *Event {
	return &Event{kind, OriginID, vectorClock, localSequenceNumber, originSequenceNumber, metadata} // new event
}

func (e *Event) HigherSequenceNumber(seqNum s.SeqNum) bool { /// come back later
	return e.LocalSequenceNumber > seqNum
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
	result.WriteString(string(e.OriginID))
	result.WriteString("\n")
	result.WriteString(e.VectorClock.String())
	result.WriteString("\n")
	result.WriteString(e.LocalSequenceNumber.String())
	result.WriteString("\n")
	result.WriteString(e.OriginSequenceNumber.String())
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
