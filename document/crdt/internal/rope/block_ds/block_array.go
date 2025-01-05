package block_ds

import (
	"fmt"
	"strings"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
)

type BlockArray struct {
	blocks   []*Block
	size     int
	realSize int // size ignoring deleted blocks
}

func NewBlockArray(b []*Block) BlockDS {
	size := 0
	realSize := 0
	blocks := make([]*Block, len(b), len(b)*2) // maybe the capacity can change.
	for i := 0; i < len(b); i++ {
		size += b[i].Len()
		if !b[i].IsDeleted() {
			realSize += b[i].Len()
		}
		blocks[i] = CopyBlock(b[i])
	}
	return &BlockArray{
		blocks,
		size,
		realSize,
	}
}

func (b *BlockArray) Size() int {
	return b.size
}
func (b *BlockArray) RealSize() int {
	return b.realSize
}
func (b *BlockArray) Len() int {
	return len(b.blocks)
}
func (b *BlockArray) Find(index int, addDeleted bool) (block *Block, localIndex int, blockIndex int) {
	if b == nil || b.blocks == nil || len(b.blocks) == 0 {
		return nil, 0, 0
	}
	if b.Get(0).Len() == 0 && index == 0 && (addDeleted || (!addDeleted && !b.Get(0).IsDeleted())) { //first block
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
			if !addDeleted {
				if block.IsDeleted() {
					continue
				}
			}
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
	removedLen := 0
	for _, blk := range blocks {
		addedLen += blk.Len()
	}
	for _, blk := range deletedBlks {
		removedLen += blk.Len()
	}
	b.blocks = append(b.blocks[:blkIndex], append(blocks, b.blocks[endOfDeletion:]...)...)

	b.updateSize(addedLen - removedLen)
	return nil
}

func (b *BlockArray) append(blocks []*Block) error {
	b.blocks = append(b.blocks, blocks...)
	addedLen := 0
	for _, blk := range blocks {
		addedLen += blk.Len()
	}
	b.updateSize(addedLen)
	return nil
}
func (b *BlockArray) String(addDeleted bool, blockSeparator string) string {
	result := strings.Builder{}
	for _, blk := range b.blocks {
		if !addDeleted {
			if !blk.IsDeleted() {
				result.WriteString(blk.Content())
				result.WriteString(blockSeparator)
			}
		} else {
			result.WriteString(blk.Content())
			result.WriteString(blockSeparator)
		}
	}
	if result.String() == "" {
		return ""
	}
	return result.String()[0 : result.Len()-len(blockSeparator)]
}
func (b *BlockArray) Split(index int, tolerance int) (BlockDS, BlockDS) { // add tolerance later for efficiency
	block, localIndex, blockIndex := b.Find(index, false) // may change
	leftLength := localIndex - 0
	rightLength := block.Len() - localIndex - 1
	if tolerance > leftLength || tolerance > rightLength {
		if rightLength < leftLength && blockIndex < b.Len() {
			blockIndex++
		}
		return NewBlockArray(b.blocks[0:blockIndex]), NewBlockArray(b.blocks[blockIndex:b.Len()])
	}

	left, right := block.Split(localIndex)
	leftBlocks := b.blocks[0:blockIndex]
	if left != nil {
		leftBlocks = append(leftBlocks, left)
	}
	rightBlocks := []*Block{}
	if blockIndex+1 < b.Len() {
		rightBlocks = b.blocks[blockIndex+1 : b.Len()]
	}
	if right != nil {
		rightBlocks = append([]*Block{right}, rightBlocks...)
	}
	return NewBlockArray(leftBlocks), NewBlockArray(rightBlocks)
}
func (b *BlockArray) Merge(blockDs BlockDS, prepend bool) BlockDS {
	index := b.Len()
	if prepend {
		index = 0
	}
	blockArray, ok := blockDs.(*BlockArray)
	if !ok {
		return nil
	}
	b.Update(index, blockArray.blocks, 0)
	return b
}
func (b *BlockArray) updateSize(diff int) { // deletion maybe needed
	b.size += diff
	b.realSize += diff
}
