package crdt_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/action"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block_ds"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/value"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

func TestPrepare(t *testing.T) {
	/// pre-test
	rope := rope.NewRope(10, 0.7, 0.6, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	rope.Insert(
		block.NewBlock(vector_clock.NewClockOffset(vector_clock.NewVectorClock("A"), 0), "ABCDE", value.ByteBuffer, false),
		vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1),
		0)

	rope.Insert(
		block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), "FGHIJ", value.ByteBuffer, false),
		vector_clock.NewClockOffset(vector_clock.VectorClock{}, 2),
		0)
	rope.PrintRope(false)

	ropeCRDT := crdt.NewCRDT(rope, "A",vector_clock.VectorClock{"A":2})

	// actions and events metadata
	actionInsertionMD, ok := action.NewInsertion("123", 6).(*action.InsertionMetadata)
	if !ok {
		fmt.Println("Error casting action insertion metadata")
	}
	eventInsertionMD, ok := event.NewInsertionEventMetadata(
		block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 3}, 0), "123", value.ByteBuffer, false),
		vector_clock.NewClockOffset(vector_clock.NewVectorClock("A"), 5),
		6,
	).(*event.InsertionEventMetadata)
	if !ok {
		fmt.Println("Error casting event insertion metadata")
	}
	///
	testCases := []struct {
		label         string
		action        *action.Action
		expectedEvent *event.Event
	}{
		{
			"Insertion action",
			action.NewAction(action.Insert, global.UserID("Samy"), global.ReplicaID("A"), actionInsertionMD),
			event.NewEvent(event.Insert, global.UserID("Samy"), global.ReplicaID("A"), vector_clock.VectorClock{"A": 3},
				eventInsertionMD),
		},
	}
	for _, ts := range testCases {
		t.Run(fmt.Sprintf("Test: %s", ts.label), func(t *testing.T) {
			want := ts.expectedEvent
			got, err := ropeCRDT.Prepare(ts.action)
			if err != nil {
				t.Errorf("unexpected error :%v", err)
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("expected \n%v, \ngot \n%v", want.String(), got.String())
			}
		})
	}
}
