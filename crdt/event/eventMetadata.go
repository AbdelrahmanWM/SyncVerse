package event

import (
	"github.com/AbdelrahmanWM/SyncVerse/crdt/global"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	"github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
)

type InsertionEventMetadata struct {
	contentBlock      *block.Block
	toBeInsertedAfter *vector_clock.ClockOffset
	startIndex int
}

func NewInsertionEventMetadata(inputs ...any) any {
	if len(inputs) == 3 {
		if contentBlock, ok := inputs[0].(*block.Block); ok {
			if toBeInsertedAfter, ok := inputs[1].(*vector_clock.ClockOffset); ok {
				if index,ok:=inputs[2].(int);ok{
					return &InsertionEventMetadata{contentBlock, toBeInsertedAfter,index}
				}
			}
		}
	}
	return nil
}

type DeletionEventMetadata struct {
	deletionMetadata global.ModifyMetadata
	startIndex int
}

func NewDeletionEventMetadata(inputs ...any) any {
	if len(inputs) == 2 {
		if deletionMetadata, ok := inputs[0].(global.ModifyMetadata); ok {
			if index,ok:=inputs[1].(int);ok{
				return &DeletionEventMetadata{deletionMetadata,index}
			}
		}
	}
	return nil
}
