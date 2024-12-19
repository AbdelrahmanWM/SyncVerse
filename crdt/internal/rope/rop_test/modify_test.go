package rope_test

import (
	"testing"

	"github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/format"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestModify(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rope := NewRope(10, 0.7, 0.6, "ropeBuffer", "blockArray", "A")
		block := block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 0}, 0), "ABC", "ropeBuffer", false)
		rope.Insert(block, vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), 0)
		want := false
		got := block.FormatExists(action.Bold)
		Assert(t, got, want)
		// adding bold effect
		rope.Modify([]ModifyMetadata{{block.ClockOffset(), [2]int{0, 3}}}, format.Format{action.Bold, ""}, 0)
		want = true
		got = block.FormatExists(action.Bold)
		Assert(t, got, want)
		// removing bold effect
		rope.Modify([]ModifyMetadata{{block.ClockOffset(), [2]int{0, 3}}}, format.Format{action.Bold, "del"}, 0)
		want = false
		got = block.FormatExists(action.Bold)
		Assert(t, got, want)

	})
}
