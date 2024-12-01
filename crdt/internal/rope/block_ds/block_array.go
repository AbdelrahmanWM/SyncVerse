package block_ds

import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
)

type BlockArray struct {
	blocks []*Block
	size   int
}

func NewBlockArray(b []*Block) BlockDS {
	size := 0
	blocks := make([]*Block, len(b), len(b)*2) // maybe the capacity can change.
	for i := 0; i < len(b); i++ {
		size += b[i].Len()
		blocks[i] = CopyBlock(b[i])
	}
	return BlockArray{
		blocks,
		size,
	}
}

func (b BlockArray) Size() int {
	return b.size
}
func (b BlockArray) Len() int {
	return len(b.blocks)
}
func (b BlockArray) Find(index int) *Block {
	if index < 0 || index >= b.size {
		return nil
	}
	for i := 0; i < len(b.blocks); i++ {
		if b.blocks[i].IsDeleted() {
			continue
		}
		block := b.blocks[i]
		if block.Len() > index {
			return block
		} else {
			index -= block.Len()
		}
	}
	return nil
}
