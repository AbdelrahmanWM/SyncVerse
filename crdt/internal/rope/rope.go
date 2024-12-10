package rope

import (
	"fmt"

	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	blockDS "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block_ds"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/node"
	vectorClock "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
	// event "github.com/AbdelrahmanWM/SyncVerse/crdt/Event"
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
	leftBlocks := make([]*Block, 1, 10)
	firstBlock := NewBlock(vectorClock.NewClockOffset(vectorClock.NewVectorClock(replicaID), 0), " ", ropeType, true)
	leftBlocks[0] = firstBlock
	root.SetLeft(NewLeafNode(blockDS.NewBlockDS(blockDSType, leftBlocks), root)) // capacity can be modified
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
func (r *Rope) Insert(contentBlock *Block, clockOffset *vectorClock.ClockOffset, startIndex int) bool {
	curNode, block, localIdx, blockIdx := r.findNodeAndBlockAndBlockIndexFromClockOffset(clockOffset, startIndex)
	fmt.Println(curNode)
	fmt.Println(block)
	fmt.Println(localIdx)
	fmt.Println(blockIdx)
	fmt.Println("****")
	/// deal with the case where its in the middle of the block
	if block.ContainsOffset(localIdx) {
		leftBlock, rightBlock := block.Split(localIdx)
		curNode.Blocks().Update(blockIdx, []*Block{leftBlock, contentBlock, rightBlock}, 1)
		return true
	}
	fmt.Println("3\n")
	var i int = blockIdx
	var node *LeafNode
	inserted := false
	nextNode := r.nextLeaf(node)
	fmt.Println("4\n")
nested:
	for node := curNode; nextNode != nil; node = nextNode {
		blk := node.Blocks().Get(i)
		for ; i < node.Blocks().Len(); i++ {

			val, err := blk.Compare(contentBlock)
			if err != nil {
				return false /// change later
			}
			switch val {

			case 0:
				if blk.CompareHashes(contentBlock) > 0 {
					node.Blocks().Update(i, []*Block{contentBlock}, 0) // insert the block in the right place
					inserted = true
					break nested

				}
			case 1:
				node.Blocks().Update(i, []*Block{contentBlock}, 0)
				inserted = true
				break nested
			}
		}
		i = 0
		nextNode = r.nextLeaf(nextNode)
	}
	fmt.Println("5\n")
	if !inserted { // inserting in the last node
		node.Blocks().Update(node.Blocks().Len(), []*Block{contentBlock}, 0)
	}
	return true
}
func (r *Rope) findNodeAndBlockAndBlockIndexFromClockOffset(clockOffset *vectorClock.ClockOffset, startIndex int) (node *LeafNode, block *Block, localIndex int, blockIndex int) { // not tested
	node, idx := r.Find(startIndex)
	block, blkLocalIdx, blockIndex := r.FindBlockFromNode(node, idx)
	if block == nil {
		return node, nil, 0, 0
	}
	for !(block.HasVectorClock(clockOffset.VectorClock()) && (block.ContainsOffset(clockOffset.Offset())) || clockOffset.Offset() == block.Len()) {
		block = node.Blocks().Get(blkLocalIdx)
		if block == nil {
			node = r.nextLeaf(node)
			if node == nil {
				return nil, nil, 0, 0 /// no blocks found
			}
			blkLocalIdx = 0
		} else {
			blkLocalIdx++
		}
	}
	return node, block, clockOffset.Offset() - block.Offset(), blkLocalIdx // get relative index within the block
}

func (r *Rope) FindBlockFromIndex(position int) (*Block, int, int) {
	node, index := r.Find(position)
	if node == nil {
		return nil, 0, 0
	}
	block, bIndex, blkIndex := node.Blocks().Find(index)
	if block == nil {
		return nil, index, blkIndex // for a later use
	}
	return block, bIndex, node.Blocks().Len()
}
func (r *Rope) FindBlockFromNode(node *LeafNode, index int) (block *Block, localIndex int, blockIndex int) {
	block, bIndex, blkIndex := node.Blocks().Find(index)
	if block == nil {
		return nil, index, node.Blocks().Len() // for a later use
	}
	return block, bIndex, blkIndex
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
