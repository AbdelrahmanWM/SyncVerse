package rope_test

import (
	"math/rand"
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

		firstBlockOffset := NewClockOffset(VectorClock{}, 0)
		startIndex := 0
		newBlock := NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "ABC", "ropeBuffer", false)
		rope.Insert(newBlock, firstBlockOffset, startIndex)
		got := rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(0)
		if !reflect.DeepEqual(got, newBlock) {
			t.Errorf("expected %v, got %v", newBlock, got)
		}
	})
	t.Run("Multiple insertions", func(t *testing.T) {
		rope := NewRope(64, 0.75, 0.65, "ropeBuffer", "blockArray", "A")

		blockOffset := NewClockOffset(VectorClock{}, 0)
		startIndex := 0
		newBlock := NewBlock(NewClockOffset(VectorClock{"B": 1}, 0), "ABC", "ropeBuffer", false)
		rope.Insert(newBlock, blockOffset, startIndex)
		got := rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(0)
		Assert(t, got, newBlock)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 1, "B": 1}, 0), "D", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{"B": 1}, 3)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Left().(*LeafNode).Blocks().(*BlockArray).Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 0, "B": 1, "C": 1}, 0), "F", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{}, 2)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Right().(*LeafNode).Blocks().Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 0, "B": 1, "C": 2}, 0), "G", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{}, 2)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Right().(*LeafNode).Blocks().Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 1, "B": 2}, 0), "H", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{}, 2)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Right().(*LeafNode).Blocks().Get(3)
		Assert(t, got, newBlock)
		rope.PrintRope(false)

		newBlock = NewBlock(NewClockOffset(VectorClock{"A": 0, "B": 0, "D": 1}, 0), "E", "ropeBuffer", false)
		blockOffset = NewClockOffset(VectorClock{}, 2)
		rope.Insert(newBlock, blockOffset, startIndex)
		got = rope.Root().Right().(*LeafNode).Blocks().Get(1)
		Assert(t, got, newBlock)
		rope.PrintRope(false)
	})
	t.Run("Inserting in different orders (concurrent Insertions)", func(t *testing.T) {
		rope := NewRope(64, 0.75, 0.65, "ropeBuffer", "blockArray", "A")
		actions := []struct {
			block       *Block
			clockOffset *ClockOffset
		}{
			{

				NewBlock(NewClockOffset(VectorClock{"B": 1}, 0), "B", "ropeBuffer", false),
				NewClockOffset(VectorClock{}, 1),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"A": 1}, 0), "A", "ropeBuffer", false),
				NewClockOffset(VectorClock{}, 1),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"C": 1}, 0), "C", "ropeBuffer", false),
				NewClockOffset(VectorClock{}, 1),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"F": 1}, 0), "F", "ropeBuffer", false),

				NewClockOffset(VectorClock{}, 2),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"G": 1}, 0), "G", "ropeBuffer", false),
				NewClockOffset(VectorClock{}, 2),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"H": 1}, 0), "H", "ropeBuffer", false),

				NewClockOffset(VectorClock{}, 2),
			},
			{

				NewBlock(NewClockOffset(VectorClock{"E": 1}, 0), "E", "ropeBuffer", false),
				NewClockOffset(VectorClock{}, 2),
			},
		}
		rand.Shuffle(len(actions), func(i, j int) {
			actions[i], actions[j] = actions[j], actions[i]
		})

		for _, action := range actions {
			rope.Insert(action.block, action.clockOffset, 0)
		}
		want := " ABC EFGH"
		got := rope.String(false)
		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}

	})
	t.Run("test insertion node split", func(t *testing.T) {
		rope := NewRope(10, 0.70, 0.65, "ropeBuffer", "blockArray", "A")
		block := NewBlock(NewClockOffset(VectorClock{"C": 1}, 0), "0123456789ABCDEF", "ropeBuffer", false)
		clock := NewClockOffset(VectorClock{}, 1)
		rope.Insert(block, clock, 0)
		// if rope.Root().Right().Right()
		rope.PrintRope(false)
		t.Fail()

	})
}

func Assert(t *testing.T, got, want *Block) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
