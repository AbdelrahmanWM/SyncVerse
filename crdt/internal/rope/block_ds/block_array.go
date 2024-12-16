package block_ds

import (
	"fmt"
	"strings"

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
	return &BlockArray{
		blocks,
		size,
	}
}

func (b *BlockArray) Size() int {
	return b.size
}
func (b *BlockArray) Len() int {
	return len(b.blocks)
}
func (b *BlockArray) Find(index int) (block *Block, localIndex int, blockIndex int) {
	if b.Get(0).Len() == 0 && index == 0 { //first block
		return b.Get(0), 0, 0
	}
	if index < 0 || index >= b.size {
		return nil, 0, 0
	}
	for i := 0; i < len(b.blocks); i++ {
		// if b.blocks[i].IsDeleted() {  // for now, the deleted contributes to the weight
		// 	continue
		// }
		block := b.blocks[i]
		if block.Len() > index {
			return block, index, i
		} else {
			index -= block.Len()
		}
	}
	return nil, 0, len(b.blocks)
}
func (b *BlockArray) Get(blockIndex int) *Block {
	if blockIndex < 0 || blockIndex >= b.Len() {
		return nil
	}
	return b.blocks[blockIndex]
}
func (b *BlockArray) NextBlock(blockIndex int) *Block {
	if blockIndex >= b.Len()-1 {
		return nil
	}
	return b.blocks[blockIndex+1]
}

func (b *BlockArray) Update(blkIndex int, blocks []*Block, numberOfDeletedBlocks int) error {
	if blkIndex < 0 {
		return fmt.Errorf("[ERROR] Invalid Index %d", blkIndex)
	} else if blkIndex >= b.Len() {
		err := b.append(blocks)
		return err
	}

	endOfDeletion := min(blkIndex+numberOfDeletedBlocks, b.Len())

	deletedBlks := b.blocks[blkIndex:endOfDeletion]

	addedLen := 0
	for _, blk := range blocks {
		addedLen += blk.Len()
	}
	for _, blk := range deletedBlks {
		addedLen -= blk.Len()
	}
	b.blocks = append(b.blocks[:blkIndex], append(blocks, b.blocks[endOfDeletion:]...)...)

	b.size += addedLen
	return nil
}

func (b *BlockArray) append(blocks []*Block) error {
	b.blocks = append(b.blocks, blocks...)
	addedLen := 0
	for _, blk := range blocks {
		addedLen += blk.Len()
	}
	b.size += addedLen
	return nil
}
func (b *BlockArray) String(addDeleted bool, blockSeparator string) string {
	result := strings.Builder{}
	for _, blk := range b.blocks {
		if addDeleted {
			if !blk.IsDeleted() {
				result.WriteString(blk.String())
				result.WriteString(blockSeparator)
			}
		} else {
			result.WriteString(blk.String())
			result.WriteString(blockSeparator)
		}
	}
	return result.String()[0 : result.Len()-len(blockSeparator)]
}
func (b *BlockArray) Split(index int) (BlockDS, BlockDS) { // add tolerance later for efficiency
	block, localIndex, blockIndex := b.Find(index)
	left, right := block.Split(localIndex)
	return NewBlockArray(append(b.blocks[0:blockIndex], left)), NewBlockArray(append([]*Block{right}, b.blocks[blockIndex+1:b.Len()]...))
}
