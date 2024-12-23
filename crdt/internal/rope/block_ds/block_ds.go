package block_ds

import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
)

type BlockDSConstructor = func(blocks []*Block) BlockDS

var BlockDSRegistry = make(map[string]BlockDSConstructor) // not optimal, change the func input later
type BlockDS interface {
	Len() int  // number of blocks
	Size() int // number of characters
	RealSize()int // excluding deleted blocks
	Find(index int, addDeleted bool) (block *Block, localIndex int, blockIndex int)
	Get(blockIndex int) *Block
	NextBlock(blockIndex int) *Block
	Update(index int, blocks []*Block, numberOfDeletedBlocks int) error
	Merge(blockds BlockDS,prepend bool)BlockDS
	String(showDeleted bool, blockSeparator string) string
	Split(index int,tolerance int)(BlockDS,BlockDS)
}

func NewBlockDS(typename string, blocks []*Block) BlockDS {
	constructor, ok := BlockDSRegistry[typename]
	if ok {
		return constructor(blocks)
	}
	return nil
}
func RegisterNewBlockDSType(typename string, constructor BlockDSConstructor) {
	BlockDSRegistry[typename] = constructor
}
func init() {
	RegisterNewBlockDSType("blockArray", NewBlockArray)
}
