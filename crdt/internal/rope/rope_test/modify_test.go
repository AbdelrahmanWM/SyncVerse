package rope_test

import (
	"testing"

	"github.com/AbdelrahmanWM/SyncVerse/crdt/action"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/global"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block_ds"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/format"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/value"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestModify(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rope := NewRope(10, 0.7, 0.6, value.ByteBuffer, block_ds.BlockArrayDS, "A")
		block := block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 0}, 0), "ABC",value.ByteBuffer, false)
		rope.Insert(block, vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), 0)
		want := false
		got := block.FormatExists(action.Bold)
		Assert(t, got, want)
		// adding bold effect
		rope.Modify([]global.ModifyMetadata{{block.ClockOffset(), [2]int{0, 3}}}, format.Format{action.Bold, ""}, 0)
		want = true
		got = block.FormatExists(action.Bold)
		Assert(t, got, want)
		// removing bold effect
		rope.Modify([]global.ModifyMetadata{{block.ClockOffset(), [2]int{0, 3}}}, format.Format{action.Bold, "del"}, 0)
		want = false
		got = block.FormatExists(action.Bold)
		Assert(t, got, want)

	})
}
