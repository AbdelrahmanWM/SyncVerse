package block_ds

import (
	"reflect"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestFind(t *testing.T) { // outdated
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
		// {
		// 	"some deleted blocks",
		// 	NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false),
		// 		NewBlock(NewClockOffset(VectorClock{"A": 2}, 0), "12345", "ropeBuffer", true),
		// 		NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "#####", "ropeBuffer", false)}),
		// 	NewBlock(NewClockOffset(VectorClock{"A": 3}, 0), "#####", "ropeBuffer", false),
		// 	5,
		// },
	}
	for _, ts := range testcases {
		t.Run(ts.description, func(t *testing.T) {
			got, _, _ := ts.target.Find(ts.index)
			if !reflect.DeepEqual(got, ts.want) {
				t.Errorf("expected %v, got %v", ts.want, got)
			}
		})
	}
}

// continue
func TestUpdate(t *testing.T) { // outdated
	testcases := []struct {
		description   string
		target        BlockDS
		index         int
		content       []*Block
		deletedNumber int
		want          BlockDS
	}{
		{
			"appending a block",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}),
			99, //append
			[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)},
			0,
			&BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}, 10},
		},
		{
			"adding in the middle",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}),
			1,
			[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "###", "ropeBuffer", false)},
			0,
			&BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "###", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}, 11},
		},
		{
			"deleting blocks from the middle",
			NewBlockArray([]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "###", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "block", "ropeBuffer", false)}),
			1, //append
			[]*Block{},
			2,
			&BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false)}, 3},
		},
	}
	for _, ts := range testcases {
		t.Run(ts.description, func(t *testing.T) {
			err := ts.target.Update(ts.index, ts.content, ts.deletedNumber)
			if err != nil {
				t.Fatalf("[ERROR] %v", err)
			}
			if !reflect.DeepEqual(ts.target, ts.want) {
				t.Errorf("expected %v, got %v", ts.want, ts.target)
			}
		})
	}
}
func TestSplit(t *testing.T) {
	testCases := []struct {
		blockArray BlockArray
		splitIndex int
		tolerance  int
		Left       string
		Right      string
	}{
		{
			BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "456", "ropeBuffer", false)}, 6},
			2,
			1,
			"123",
			"456",
		},
		{
			BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "456", "ropeBuffer", false)}, 6},
			2,
			0,
			"12",
			"3456",
		},
		{
			BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "456", "ropeBuffer", false)}, 6},
			4,
			3,
			"123",
			"456",
		},
		{
			BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "123", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "456", "ropeBuffer", false)}, 6},
			5,
			2,
			"123456",
			"",
		},
		{
			BlockArray{[]*Block{NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "12345", "ropeBuffer", false), NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "6789A", "ropeBuffer", false)}, 10},
			7,
			2,
			"1234567",
			"89A",
		},
	}
	for i, ts := range testCases {
		want1, want2 := ts.Left, ts.Right
		got1B, got2B := ts.blockArray.Split(ts.splitIndex, ts.tolerance)
		got1, got2 := got1B.String(false, ""), got2B.String(false, "")
		if !reflect.DeepEqual(want1, got1) {
			t.Errorf("%d [LEFT] expected %v, got %v", i, want1, got1)
		}
		if !reflect.DeepEqual(want2, got2) {
			t.Errorf("[Right] expected %v, got %v", want2, got2)
		}
	}
}
