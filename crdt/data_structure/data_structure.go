package data_structure

import (
	"github.com/AbdelrahmanWM/SyncVerse/crdt/global"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type CRDTDataStructure interface {
	Insert(contentBlock *block.Block, toBeInsertedAfterOffset *vector_clock.ClockOffset, startIndex int)error
	Delete(toBeDeletedBlocksMetadata []global.ModifyMetadata, startIndex int )error
	FindBlocks(firstBlockStartIndex int, length int)([]global.ModifyMetadata, error)
	FindInsertionBlockOffset(insertionPosition int)(clockOffset *vector_clock.ClockOffset,err error) 
}
