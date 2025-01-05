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

	// actions and events metadata
	actionInsertionMD, ok := action.NewInsertion("123", 6).(*action.InsertionMetadata)
	if !ok {
		t.Error("Error casting action insertion metadata")
	}
	eventInsertionMD, ok := event.NewInsertionEventMetadata(
		block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 3}, 0), "123", value.ByteBuffer, false),
		vector_clock.NewClockOffset(vector_clock.NewVectorClock("A"), 5),
		6,
	).(*event.InsertionEventMetadata)
	if !ok {
		t.Error("Error casting event insertion metadata")
	}
	////////////////////////////////////
	actionDeletionMD, ok := action.NewDeletion(6, 0).(*action.DeletionMetadata)
	if !ok {
		t.Error("Error casting deletion metadata")
	}
	eventDeletionMD, ok := event.NewDeletionEventMetadata(global.ModifyMetadataArray{&global.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 0), [2]int{0, 1}}, &global.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 1}, 0), [2]int{0, 5}}}, 0).(*event.DeletionEventMetadata)
	if !ok {
		t.Error("Error casting deletion event metadata ")
	}
	///
	testCases := []struct { // independent
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
		{
			"Deletion action",
			action.NewAction(action.Delete, global.UserID("Hany"), global.ReplicaID("A"), actionDeletionMD),
			event.NewEvent(event.Delete, global.UserID("Hany"), global.ReplicaID("A"), vector_clock.VectorClock{"A": 3}, eventDeletionMD),
		},
	}
	for i, ts := range testCases {
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

		ropeCRDT := crdt.NewCRDT(rope, "A", vector_clock.VectorClock{"A": 2})

		// test
		t.Run(fmt.Sprintf("%d- %s", i, ts.label), func(t *testing.T) {
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

func TestApply(t *testing.T) {
	// events metadata
	event1MD, ok := event.NewInsertionEventMetadata(
		block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 3}, 0), "123", value.ByteBuffer, false),
		vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 1}, 5),
		5).(*event.InsertionEventMetadata)
	if !ok {
		t.Error("Error casting event1 metadata")
	}
    event2MD,ok:=event.NewDeletionEventMetadata(
		global.ModifyMetadataArray{&global.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{},1),[2]int{0,1}},&global.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"A":2},0),[2]int{0,5}},},
		5,
	).(*event.DeletionEventMetadata)
	if !ok{
		t.Error("Error casting event2 metadata")
	}
	testcases := []struct { // independent 
		label                        string
		event                        *event.Event
		expectedRopeStringAfterEvent string
	}{
		{
			"Insertion event",
		event.NewEvent(event.Insert,global.UserID("Samy"),global.ReplicaID("A"),vector_clock.VectorClock{"A":3},event1MD),
		" ABCDE123 FGHIJ",
		},
		{
			"Deletion event",
			event.NewEvent(event.Delete,global.UserID("Samy"),global.ReplicaID("A"),vector_clock.VectorClock{"A":3},event2MD),
			" ABCDE",
		},
	}
	for i, tc := range testcases {
		// pretest
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
		fmt.Println(rope.String(false))

		ropeCRDT := crdt.NewCRDT(rope, "A", vector_clock.VectorClock{"A": 2})
		t.Run(fmt.Sprintf("%d- %s", i, tc.label), func(t *testing.T) {
			want := tc.expectedRopeStringAfterEvent
			err := ropeCRDT.Apply(tc.event)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			addDeletedBlocks := false
			got := ropeCRDT.DataStructure().String(addDeletedBlocks)
			if want != got {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}
