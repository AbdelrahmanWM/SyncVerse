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
	root := NewInnerNode(1, 2, nil, nil, nil)
	leftBlocks := make([]*Block, 1, 10)
	firstVectorClock := vectorClock.NewVectorClock(replicaID)
	leftBlocks[0] = NewBlock(vectorClock.NewClockOffset(firstVectorClock, 0), " ", ropeType, false)
	rightBlocks := make([]*Block, 1, 10)
	rightBlocks[0] = NewBlock(vectorClock.NewClockOffset(firstVectorClock.NewVectorClock(replicaID), 1), " ", ropeType, false)
	root.SetLeft(NewLeafNode(blockDS.NewBlockDS(blockDSType, leftBlocks), root)) // capacity can be modified
	root.SetRight(NewLeafNode(blockDS.NewBlockDS(blockDSType, rightBlocks), root))

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
	inserted := false

	/// deal with the case where its in the middle of the block
	if block.ContainsOffset(localIdx) {
		leftBlock, rightBlock := block.Split(localIdx)
		blocks := []*Block{}
		if leftBlock != nil {
			blocks = append(blocks, leftBlock)
		}
		blocks = append(blocks, contentBlock)
		blocks = append(blocks, rightBlock)
		curNode.Blocks().Update(blockIdx, blocks, 1)
		inserted = true
	}
	var i int = blockIdx + 1
	var node LeafNode = *curNode
	var refNode *LeafNode = &node
	var len int
	nextNode := r.nextLeaf(curNode)

	for node = *curNode; nextNode != nil && !inserted; node, nextNode = *nextNode, r.nextLeaf(nextNode) {
		refNode = &node
		blk := refNode.Blocks().Get(i)
		len = refNode.Blocks().Len()
		for ; i < len; i++ {

			val, err := blk.Compare(contentBlock)
			if err != nil {
				return false /// change later
			}
			switch val {

			case 0: // concurrent
				hashComp := blk.CompareHashes(contentBlock)
				if hashComp > 0 {
					refNode.Blocks().Update(i, []*Block{contentBlock}, 0) // insert the block in the right place
					inserted = true
				} else if hashComp == 0 {
					fmt.Println("[ERROR] SAME EVENT VECTOR CLOCK!")
					return false
				} else {
					if i == len-1 {

						nextBlock := nextNode.Blocks().Get(0)
						compare, err := nextBlock.Compare(contentBlock)
						if err != nil {
							return false
						}
						if compare == -1 || (compare == 0 && nextBlock.CompareHashes(contentBlock) > 0) {
							inserted = true
						}

					}
				}
			case -1: // already known event
				refNode.Blocks().Update(i, []*Block{contentBlock}, 0)
				inserted = true
			case 1: // gap in events
				fmt.Println("[ERROR] NON-POSSIBLE EVENT ORDER")

				return false
			}
			if inserted {
				break
			}
		}
		i = 0

	}

	if !inserted { // inserting in the last node
		refNode.Blocks().Update(refNode.Blocks().Len(), []*Block{contentBlock}, 0)
	}
	updateWeight(refNode, contentBlock.Len())
	return true
}
func updateWeight(node RopeNode, diff int) {

	if node == nil || node.Parent() == nil {
		return
	}
	parent := node.Parent().(*InnerNode)
	if isLeftChild(node) {
		parent.SetLeftWeight(parent.LeftWeight() + diff)
	}
	parent.SetWeight(parent.Weight() + diff)
	updateWeight(parent, diff)
}
func (r *Rope) findNodeAndBlockAndBlockIndexFromClockOffset(clockOffset *vectorClock.ClockOffset, startIndex int) (node *LeafNode, block *Block, localIndex int, blockIndex int) { // not tested
	node, idx := r.Find(startIndex)
	if node == nil {
		return nil, nil, 0, 0
	}
	block, blkLocalIdx, blockIndex := r.FindBlockFromNode(node, idx)
	if block == nil {
		return node, nil, 0, 0
	}
	i := blockIndex
	len := node.Blocks().Len()
	for node != nil {
		for ; i < len; i++ {
			if block.HasVectorClock(clockOffset.VectorClock()) && (block.ContainsOffset(clockOffset.Offset()) || block.Len() == clockOffset.Offset()) {
				return node, block, clockOffset.Offset() - block.Offset(), blkLocalIdx // get relative index within the block
			}
		}
		node = r.nextLeaf(node)
		i = 0
		len = node.Blocks().Len()
	}
	return nil, nil, 0, 0

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
func (r *Rope) PrintRope(addDeleted bool) {
	for queue := []RopeNode{r.Root()}; len(queue) > 0; {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			switch castedNode := node.(type) {
			case *LeafNode:
				fmt.Printf(" %v|<%v> ", castedNode.Weight(), castedNode.Blocks().String(addDeleted))
			case *InnerNode:
				fmt.Printf(" %v|%v ", castedNode.LeftWeight(), castedNode.Weight())

			}
			queue = queue[1:]
			if node.Left() != nil {
				queue = append(queue, node.Left())
			}
			if node.Right() != nil {
				queue = append(queue, node.Right())
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
