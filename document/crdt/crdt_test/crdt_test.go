package crdt_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/action"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/event"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block_ds"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/value"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/sequence_number"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/types"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/document/global"
)

const (
	IGNORED = 0
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
	eventDeletionMD, ok := event.NewDeletionEventMetadata(types.ModifyMetadataArray{&types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 0), [2]int{0, 1}}, &types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 1}, 0), [2]int{0, 5}}}, 0).(*event.DeletionEventMetadata)
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
			action.NewAction(action.Insert, actionInsertionMD),
			event.NewEvent(event.Insert, global.ReplicaID("A"), vector_clock.VectorClock{"A": 3}, 1, 1,
				eventInsertionMD),
		},
		{
			"Deletion action",
			action.NewAction(action.Delete, actionDeletionMD),
			event.NewEvent(event.Delete, global.ReplicaID("A"), vector_clock.VectorClock{"A": 3}, 1, 1, eventDeletionMD),
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
	t.Run("Limited event application effect comparison test", func(t *testing.T) {
		// events metadata
		event1MD, ok := event.NewInsertionEventMetadata(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 3}, 0), "123", value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 1}, 5),
			5).(*event.InsertionEventMetadata)
		if !ok {
			t.Error("Error casting event1 metadata")
		}
		event2MD, ok := event.NewDeletionEventMetadata(
			types.ModifyMetadataArray{&types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}, &types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), [2]int{0, 5}}},
			5,
		).(*event.DeletionEventMetadata)
		if !ok {
			t.Error("Error casting event2 metadata")
		}
		testcases := []struct { // independent
			label                        string
			event                        *event.Event
			expectedRopeStringAfterEvent string
		}{
			{
				"Insertion event",
				event.NewEvent(event.Insert, global.ReplicaID("A"), vector_clock.VectorClock{"A": 3}, IGNORED, IGNORED, event1MD),
				" ABCDE123 FGHIJ",
			},
			{
				"Deletion event",
				event.NewEvent(event.Delete, global.ReplicaID("A"), vector_clock.VectorClock{"A": 3}, IGNORED, IGNORED, event2MD),
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
	})


	t.Run("Test sequence numbers update when an event is applied",func(t *testing.T){

		// events metadata
		event1MD, ok := event.NewInsertionEventMetadata(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"B": 1}, 0), "123", value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, 0),
			0).(*event.InsertionEventMetadata)

		if !ok {
			t.Error("Error casting event1 metadata")
		}
		event2MD, ok := event.NewDeletionEventMetadata(
			types.ModifyMetadataArray{&types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}},
			5,
		).(*event.DeletionEventMetadata)
		if !ok {
			t.Error("Error casting event2 metadata")
		}
		event3MD, ok := event.NewInsertionEventMetadata(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"C": 1}, 0), "123", value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, 0),
			0).(*event.InsertionEventMetadata)

		if !ok {
			t.Error("Error casting event3 metadata")
		}
		event4MD, ok := event.NewDeletionEventMetadata(
			types.ModifyMetadataArray{&types.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"C":1}, 0), [2]int{0, 3}}},
			0,
			).(*event.DeletionEventMetadata)

		if !ok {
			t.Error("Error casting event4 metadata")
		}
		testcases := []struct { // sequence of dependent events
			label                        string
			event                        *event.Event
			seqNumber sequence_number.SeqNum
			seqNumbrMap sequence_number.SeqNumMap
		}{
			{
				"Sequence number map basics",
				event.NewEvent(event.Insert, global.ReplicaID("B"), vector_clock.VectorClock{"B": 1}, IGNORED, 1, event1MD),
				sequence_number.SeqNum(1),
				sequence_number.SeqNumMap{"B":1},
			},
			{
				"Sequence number map doesn't track self",
				event.NewEvent(event.Delete, global.ReplicaID("A"), vector_clock.VectorClock{"A": 1,"B":1}, IGNORED, 2, event2MD),
				sequence_number.SeqNum(2),
				sequence_number.SeqNumMap{"B":1},
			},
			{
				"event from replica C",
				event.NewEvent(event.Insert, global.ReplicaID("C"),vector_clock.VectorClock{"C":1},IGNORED,1,event3MD),
				sequence_number.SeqNum(3),
				sequence_number.SeqNumMap{"B":1,"C":1},
			},
			{
				"another event from replica C",
				event.NewEvent(event.Delete, global.ReplicaID("C"),vector_clock.VectorClock{"C":2},IGNORED,2,event4MD),
				sequence_number.SeqNum(4),
				sequence_number.SeqNumMap{"B":1,"C":2},
			},
		}

		// pretest

		rope := rope.NewRope(10, 0.7, 0.6, value.ByteBuffer, block_ds.BlockArrayDS, "A")
		ropeCRDT := crdt.NewCRDT(rope, "A", vector_clock.VectorClock{"A": 0})

		for i, tc := range testcases {
			
		t.Run(fmt.Sprintf("%d- %s", i, tc.label), func(t *testing.T) {
				seqNum := tc.seqNumber
				err := ropeCRDT.Apply(tc.event)
				if err != nil {
					t.Errorf("unexpected error %v", err)
				}
				gottenSeqNum :=ropeCRDT.LocalSequenceNumber()
				if seqNum != gottenSeqNum {
					t.Errorf("expected %v, got %v", seqNum, gottenSeqNum)
				}
				seqNumMap := tc.seqNumbrMap
				gottenSeqNumMap := ropeCRDT.ObservedSequenceNumberMap()
				if !reflect.DeepEqual(seqNumMap, gottenSeqNumMap){
					t.Errorf("expected %v, got %v", seqNumMap, gottenSeqNumMap)
				}

				
			})
		}
	})
}
