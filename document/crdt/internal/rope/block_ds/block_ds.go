package block_ds

import (
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
)

type BlockDSConstructor = func(blocks []*Block) BlockDS
type BlockDSType int

const (
	BlockArrayDS BlockDSType = iota
)

var BlockDSRegistry = make(map[BlockDSType]BlockDSConstructor) // not optimal, change the func input later
type BlockDS interface {
	Len() int      // number of blocks
	Size() int     // number of characters
	RealSize() int // excluding deleted blocks
	Find(index int, addDeleted bool) (block *Block, localIndex int, blockIndex int)
	Get(blockIndex int) *Block
	NextBlock(blockIndex int) *Block
	Update(index int, blocks []*Block, numberOfDeletedBlocks int) error
	Merge(blockds BlockDS, prepend bool) BlockDS
	String(showDeleted bool, blockSeparator string) string
	Split(index int, tolerance int) (BlockDS, BlockDS)
}

func NewBlockDS(typename BlockDSType, blocks []*Block) BlockDS {
	constructor, ok := BlockDSRegistry[typename]
	if ok {
		return constructor(blocks)
	}
	return nil
}
func RegisterNewBlockDSType(typename BlockDSType, constructor BlockDSConstructor) {
	BlockDSRegistry[typename] = constructor
}
func init() {
	RegisterNewBlockDSType(BlockArrayDS, NewBlockArray)
}
