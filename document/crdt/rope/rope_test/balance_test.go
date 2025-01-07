package rope_test

import (
	"fmt"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block_ds"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/value"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
)

func TestBalanceLeaves(t *testing.T) {  /// Failing test for examination purposes
	rope := NewRope(12, 0.7, 0.6, value.ByteBuffer, block_ds.BlockArrayDS, "A")

	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "0", value.ByteBuffer, false), NewClockOffset(VectorClock{}, 1), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "123", value.ByteBuffer, false), NewClockOffset(VectorClock{}, 2), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "45678", value.ByteBuffer, false), NewClockOffset(VectorClock{}, 1), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 4}, 0), "ABCDEFGHIJKLM", value.ByteBuffer, false), NewClockOffset(VectorClock{}, 2), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 5}, 0), "IIIII", value.ByteBuffer, false), NewClockOffset(VectorClock{}, 1), 0)
	fmt.Println("Initial Rope State:")
	rope.PrintRope(false)
	rope.BalanceLeaves(Right)
	fmt.Println("Rope state after balancing to the Right:")
	rope.PrintRope(false)
	rope.BalanceLeaves(Left)
	fmt.Println("Rope state after balancing to the Left:")
	rope.PrintRope(false)
	t.Fail()

}
