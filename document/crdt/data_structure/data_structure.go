package data_structure

import (
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/types"
)

type CRDTDataStructure interface {
	Insert(contentBlock *block.Block, toBeInsertedAfterOffset *vector_clock.ClockOffset, startIndex int) error
	Delete(toBeDeletedBlocksMetadata types.ModifyMetadataArray, startIndex int) error
	FindBlocks(firstBlockStartIndex int, length int) (types.ModifyMetadataArray, error)
	FindInsertionBlockOffset(insertionPosition int) (clockOffset *vector_clock.ClockOffset, err error)
	String(addDeletedBlocks bool) string
}
