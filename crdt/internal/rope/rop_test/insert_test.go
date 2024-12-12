package rope_test

import (
	"reflect"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block_ds"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/node"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

func TestNewRope(t *testing.T) {

	rope := NewRope(64, 0.75, 0.65, "ropeBuffer", "blockArray", "A")
	leftNode, ok := rope.Root().Left().(*LeafNode)
	if !ok {
		t.Fatal("Invalid rope structure.")
	}
	firstBlock := leftNode.Blocks().(*BlockArray).Get(0)
	if firstBlock == nil {
		t.Errorf("First block not found")
	}
}

func TestInsert(t *testing.T) {
	t.Run("Limited Insert testing", func(t *testing.T) {
		rope := NewRope(64, 0.75, 0.65, "ropeBuffer", "blockArray", "A")

		firstBlockOffset := NewClockOffset(VectorClock{"A": 1}, 0)
		startIndex := 0
		newBlock := NewBlock(NewClockOffset(VectorClock{"A": 1, "B": 1}, 0), "ABC", "ropeBuffer", false)
		rope.Insert(newBlock, firstBlockOffset, startIndex)
		got := rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(0)
		if !reflect.DeepEqual(got, newBlock) {
			t.Errorf("expected %v, got %v", newBlock, got)
		}
	})
	t.Run("Multiple insertions", func(t *testing.T) {
		rope := NewRope(64, 0.75, 0.65, "ropeBuffer", "blockArray", "A")

		blockOffset := NewClockOffset(VectorClock{"A": 1}, 0)
		startIndex := 0
		newBlock := NewBlock(NewClockOffset(VectorClock{"B": 1}, 0), "ABC", "ropeBuffer", false)
		rope.Insert(newBlock, blockOffset, startIndex)
		got := rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(0)
		Assert(t, got, newBlock)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 1, "B": 2}, 0), "D", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{"B": 1}, 3)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 1, "B": 2, "C": 1}, 0), "F", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{"A": 1}, 2)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Right().(*LeafNode).Blocks().Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)
	})
}
func Assert(t *testing.T, got, want *Block) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
