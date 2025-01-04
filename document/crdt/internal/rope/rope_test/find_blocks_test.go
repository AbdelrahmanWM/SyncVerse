package rope_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block_ds"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/value"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

func TestFindBlocks(t *testing.T) {
	rope1 := rope.NewRope(10, 0.8, 0.65, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	rope2 := rope.NewRope(5, 0.8, 0.65, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	rope3 := rope.NewRope(4, 0.75, 0.65, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	rope4 := rope.NewRope(4, 0.75, 0.65, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	rope5 := rope.NewRope(4, 0.75, 0.65, value.ByteBuffer, block_ds.BlockArrayDS, "A")
	for i := 1; i <= 10; i++ {
		rope1.Insert(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0),
				fmt.Sprintf("%d", i-1), value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, (i%2)+1), 0)
		rope2.Insert(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0),
				fmt.Sprintf("%d", i-1), value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, (i%2)+1), 0)
		rope3.Insert(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0),
				fmt.Sprintf("%d", i-1), value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, (i%2)+1), 0)
		rope4.Insert(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0),
				fmt.Sprintf("%d", i-1), value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, (i%2)+1), 0)
		rope5.Insert(
			block.NewBlock(vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0),
				fmt.Sprintf("%d", i-1), value.ByteBuffer, false),
			vector_clock.NewClockOffset(vector_clock.VectorClock{}, (i%2)+1), 0)
	}
	rope4.Delete([]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), [2]int{0, 1}}}, 0)
	rope5Deletions := []global.ModifyMetadata{}
	for i := 10; i > 0; i -= 2 {
		rope5Deletions = append(rope5Deletions, global.ModifyMetadata{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": i}, 0), [2]int{0, 1}})
	}
	rope5.Delete(rope5Deletions, 0)
	rope1.PrintRope(false)
	rope2.PrintRope(false)
	rope3.PrintRope(false)
	rope4.PrintRope(false)
	rope5.PrintRope(false)
	testCases := []struct {
		rope                 *rope.Rope
		index                int
		length               int
		modificationMetadata []global.ModifyMetadata
	}{
		{
			rope1,
			5,
			2,
			[]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), [2]int{0, 1}}, {vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}},
		},
		{
			rope2,
			5,
			2,
			[]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), [2]int{0, 1}}, {vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}},
		},
		{
			rope3,
			5,
			2,
			[]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 2}, 0), [2]int{0, 1}}, {vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}},
		},
		{
			rope4,
			5,
			2,
			[]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}, {vector_clock.NewClockOffset(vector_clock.VectorClock{"A": 9}, 0), [2]int{0, 1}}},
		},
		{
			rope5,
			0,
			2,
			[]global.ModifyMetadata{{vector_clock.NewClockOffset(vector_clock.VectorClock{}, 0), [2]int{0, 1}}, {vector_clock.NewClockOffset(vector_clock.VectorClock{}, 1), [2]int{0, 1}}},
		},
	}
	for _, ts := range testCases {
		want := ts.modificationMetadata
		got, err := ts.rope.FindBlocks(ts.index, ts.length)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		if len(want) != len(got) {
			t.Errorf("expected length of %d, got %d", len(want), len(got))
		}
		for i, _ := range want {
			if !got[i].ClockOffset.Equals(want[i].ClockOffset) {
				t.Errorf("expected %s, got %s", want[i].ClockOffset.String(), got[i].ClockOffset.String())
			} else if !reflect.DeepEqual(got[i].Rng, want[i].Rng) {
				t.Errorf("expected %v, got %v", want[i].Rng, got[i].Rng)
			}
		}
	}
}
