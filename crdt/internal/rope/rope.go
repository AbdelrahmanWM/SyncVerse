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
func (r *Rope) FindBlock(position int) (*Block, int) {
	node, index := r.Find(position)
	if node == nil {
		return nil, 0
	}
	block, bIndex := node.Blocks().Find(index)
	if block == nil {
		return nil, index // for a later use
	}
	return block, bIndex
}
func (r *Rope) FindBlockFromNode(node *LeafNode, index int)(*Block, int){
block, bIndex := node.Blocks().Find(index)
	if block == nil {
		return nil, index // for a later use
	}
	return block, bIndex
}



func (r *Rope) nextLeaf(leafNode *LeafNode) *LeafNode {
	if leafNode == nil {
		return nil
	}

	var ptr RopeNode
	ptr = leafNode

	for isRightChild(ptr) {
		ptr = ptr.Parent()
	}
	if ptr == nil {
		return nil
	}
	ptr = ptr.Parent()
	if ptr == nil {
		return nil
	}
	ptr = ptr.Right()
	for _, ok := ptr.(*LeafNode); ptr != nil && !ok; _, ok = ptr.(*LeafNode) {
		ptr = ptr.Left()
	}
	return ptr.(*LeafNode)
}
func isRightChild(childNode RopeNode) bool {
	parent := childNode.Parent()
	if parent == nil {
		return false
	}
	return childNode.Parent().Right() == childNode
}

func isLeftChild(childNode RopeNode) bool {
	parent := childNode.Parent()
	if parent == nil {
		return false
	}
	return childNode.Parent().Left() == childNode
}
func replaceChild(parent, currentChild, newChild RopeNode) {
	if parent.Left() == currentChild {
		parent.SetLeft(newChild)
	} else {
		parent.SetRight(newChild)
	}
}