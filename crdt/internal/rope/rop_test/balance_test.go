package rope_test

import (
	"fmt"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestBalanceLeaves(t *testing.T) {
	rope := NewRope(12, 0.7, 0.6, "ropeBuffer", "blockArray", "A")

	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 0}, 0), "0", "ropeBuffer", false), NewClockOffset(VectorClock{}, 1), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewClockOffset(VectorClock{}, 2), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "45678", "ropeBuffer", false), NewClockOffset(VectorClock{}, 1), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "ABCDEFGHIJKLM", "ropeBuffer", false), NewClockOffset(VectorClock{}, 2), 0)
	rope.Insert(NewBlock(NewClockOffset(VectorClock{"A": 4}, 0), "IIIII", "ropeBuffer", false), NewClockOffset(VectorClock{}, 1), 0)
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
