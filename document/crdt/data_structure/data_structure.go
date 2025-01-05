package data_structure

import (
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

type CRDTDataStructure interface {
	Insert(contentBlock *block.Block, toBeInsertedAfterOffset *vector_clock.ClockOffset, startIndex int) error
	Delete(toBeDeletedBlocksMetadata global.ModifyMetadataArray, startIndex int) error
	FindBlocks(firstBlockStartIndex int, length int) (global.ModifyMetadataArray, error)
	FindInsertionBlockOffset(insertionPosition int) (clockOffset *vector_clock.ClockOffset, err error)
	String(addDeletedBlocks bool) string
}
