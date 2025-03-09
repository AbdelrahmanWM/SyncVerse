package event

import (
	"strconv"
	"strings"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/types"
)

type EventMetadata interface {
	String() string
}
type InsertionEventMetadata struct {
	ContentBlock      *block.Block
	ToBeInsertedAfter *vector_clock.ClockOffset
	StartIndex        int
}

// (contentBlock *block.Block, toBeInsertedAfter *vector_clock.ClockOffset, startIndex int)
func NewInsertionEventMetadata(inputs ...any) EventMetadata {
	if len(inputs) == 3 {
		if contentBlock, ok := inputs[0].(*block.Block); ok {
			if toBeInsertedAfter, ok := inputs[1].(*vector_clock.ClockOffset); ok {
				if index, ok := inputs[2].(int); ok {
					return &InsertionEventMetadata{contentBlock, toBeInsertedAfter, index}
				}
			}
		}
	}
	return nil
}
func (ivm *InsertionEventMetadata) String() string {
	var result strings.Builder
	result.WriteString(ivm.ContentBlock.String())
	result.WriteString("\n")
	result.WriteString(ivm.ToBeInsertedAfter.String())
	result.WriteString("\n")
	result.WriteString("Start Index: ")
	result.WriteString(strconv.Itoa(ivm.StartIndex))
	return result.String()
}

type DeletionEventMetadata struct {
	DeletionMetadata types.ModifyMetadataArray
	StartIndex       int
}

// (deletionMetadata types.ModifyMetadataArray, startIndex int)
func NewDeletionEventMetadata(inputs ...any) EventMetadata {
	if len(inputs) == 2 {
		if deletionMetadata, ok := inputs[0].(types.ModifyMetadataArray); ok {
			if index, ok := inputs[1].(int); ok {
				return &DeletionEventMetadata{deletionMetadata, index}
			}
		}
	}
	return nil
}
func (dem *DeletionEventMetadata) String() string {
	var result strings.Builder
	result.WriteString(dem.DeletionMetadata.String())
	result.WriteString("\n")
	result.WriteString("Start Index: ")
	result.WriteString(strconv.Itoa(dem.StartIndex))
	return result.String()
}
