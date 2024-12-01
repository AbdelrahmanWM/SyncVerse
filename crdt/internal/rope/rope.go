package rope

import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	blockDS "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block_ds"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/node"
)

type Rope struct {
	root        *InnerNode
	chunkSize   int
	splitSize   int
	mergeSize   int
	ropeType    string
	blockDSType string
	size        int
	replicaID   string
}

/*
**
Initializes a new rope with two empty leaf nodes
*/
func NewRope(maximumChunkLength int, splitRatio float64, mergeRatio float64, ropeType string, blockDSType string, replicaID string) *Rope {
	root := NewInnerNode(0, 0, nil, nil, nil)
	root.SetLeft(NewLeafNode(blockDS.NewBlockDS(blockDSType, make([]*Block, 0, 10)), root)) // capacity can be modified
	root.SetRight(NewLeafNode(blockDS.NewBlockDS(blockDSType, make([]*Block, 0, 10)), root))
	splitSize := max(1, int(splitRatio*float64(maximumChunkLength)))
	mergeSize := max(1, int(mergeRatio*float64(maximumChunkLength)))

	return &Rope{
		root,
		maximumChunkLength,
		splitSize,
		mergeSize,
		ropeType,
		blockDSType,
		0,
		replicaID,
	}
}
func (r *Rope) NewRopeBlockDS(blocks []*Block) blockDS.BlockDS {
	return blockDS.NewBlockDS(r.blockDSType, blocks)
}
func (r *Rope) ReplicaID() string {
	return r.replicaID
}
func (r *Rope) Root() *InnerNode {
	return r.root
}

func (r *Rope) SetRoot(newRoot *InnerNode) {
	r.root = newRoot
}
func (r *Rope) ChunkSize() int {
	return r.chunkSize
}
func (r *Rope) Find(position int) (*LeafNode, int) {
	var ptr RopeNode
	ptr = r.Root()
	if ptr.Weight() == 0 && position == 0 { // special case 0-index weight problem
		return ptr.Left().(*LeafNode), 0
	}
	for ptr != nil {
		switch p := ptr.(type) {
		case *InnerNode:
			if position >= p.LeftWeight() { // the equal is due to the 0-index
				ptr = ptr.Right()
				position -= p.LeftWeight()
			} else {
				ptr = ptr.Left()
			}
		default:
			return ptr.(*LeafNode), position
		}
	}
	return nil, -1
}
