package block_ds

import (
	"reflect"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestFind(t *testing.T) {
	testcases := []struct {
		description string
		target      BlockDS
		want        *Block
		index       int
	}{
		{
			"happy path",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}),
			NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false),
			0,
		},
		{
			"empty",
			NewBlockArray([]*Block{}),
			nil,
			0,
		},
		{
			"multiple blocks",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false),
				NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "12345", "ropeBuffer", false),
				NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "#####", "ropeBuffer", false)}),
			NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "12345", "ropeBuffer", false),
			5,
		},
		{
			"some deleted blocks",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false),
				NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "12345", "ropeBuffer", true),
				NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "#####", "ropeBuffer", false)}),
			NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "#####", "ropeBuffer", false),
			5,
		},
	}
	for _, ts := range testcases {
		t.Run(ts.description, func(t *testing.T) {
			got,_ := ts.target.Find(ts.index)
			if !reflect.DeepEqual(got, ts.want) {
				t.Errorf("expected %v, got %v", ts.want, got)
			}
		})
	}
}
