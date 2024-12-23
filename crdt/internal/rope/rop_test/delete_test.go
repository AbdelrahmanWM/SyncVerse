package rope_test

import (
	"fmt"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"

	// . "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block_ds"
	// . "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/node"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestDelete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rope := NewRope(10, 0.70, 0.65, "ropeBuffer", "blockArray", "A")
		for i := 0; i < 10; i++ {
			rope.Insert(
				NewBlock(NewClockOffset(VectorClock{"A": i}, 0),
					fmt.Sprintf("%d", i), "ropeBuffer", false),
				NewClockOffset(VectorClock{}, (i%2)+1), 0)
		}
		want := " 86420 97531"
		got := rope.String(false)
		Assert(t, got, want)
		toBeDeleted := []ModifyMetadata{
			{
				NewClockOffset(VectorClock{"A": 0}, 0),
				[2]int{0, 1},
			},
		}
		rope.Delete(toBeDeleted, 0)
		want = " 8642 97531"
		got = rope.String(false)

		Assert(t, got, want)
	})
	t.Run("Multiple deletions", func(t *testing.T) {
		rope := NewRope(10, 0.70, 0.65, "ropeBuffer", "blockArray", "A")
		for i := 0; i < 10; i++ {
			rope.Insert(
				NewBlock(NewClockOffset(VectorClock{"A": i}, 0),
					fmt.Sprintf("%d", i), "ropeBuffer", false),
				NewClockOffset(VectorClock{}, (i%2)+1), 0)
		}
		want := " 86420 97531"
		got := rope.String(false)
		Assert(t, got, want)
		toBeDeleted := []ModifyMetadata{ // blocks have to be ordered (blocks order is fixed)
			{
				NewClockOffset(VectorClock{"A": 2}, 0),
				[2]int{0, 1},
			},
			{
				NewClockOffset(VectorClock{"A": 0}, 0),
				[2]int{0, 1},
			},
			{
				NewClockOffset(VectorClock{"A": 7}, 0),
				[2]int{0, 1},
			},
			{
				NewClockOffset(VectorClock{"A": 1}, 0),
				[2]int{0, 1},
			},
		}
		rope.Delete(toBeDeleted, 0)
		want = " 864 953"
		got = rope.String(false)

		Assert(t, got, want)
	})
	t.Run("delete part of a block", func(t *testing.T) {

		rope := NewRope(50, 0.70, 0.65, "ropeBuffer", "blockArray", "A")
		rope.Insert(
			NewBlock(NewClockOffset(VectorClock{"A": 0}, 0),
				"This is Not Deleted", "ropeBuffer", false),
			NewClockOffset(VectorClock{}, 1), 0)
		toBeDeleted := []ModifyMetadata{ // blocks have to be ordered (blocks order is fixed)
			{
				NewClockOffset(VectorClock{"A": 0}, 0),
				[2]int{8, 12},
			},
		}
		rope.Delete(toBeDeleted, 0)
		rope.PrintRope(false)
		want := " This is Deleted "
		got := rope.String(false)
		Assert(t, got, want)
	})
	t.Run("delete a divided block", func(t *testing.T) {
		rope := NewRope(50, 0.7, 0.65, "ropeBuffer", "blockArray", "A")
		rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 0}, 0), "|", "ropeBuffer", false), NewClockOffset(VectorClock{}, 1), 0)
		rope.Delete([]ModifyMetadata{{NewClockOffset(VectorClock{}, 0), [2]int{0, 2}}}, 0)
		want := "|"
		got := rope.String(false)
		rope.PrintRope(false)
		Assert(t, got, want)
	})

}
