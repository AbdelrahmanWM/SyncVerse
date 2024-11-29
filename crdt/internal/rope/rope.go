package rope
import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/node"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/vector_clock"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/value"

)

type Rope struct {
	root *InnerNode
	chunkSize int
	splitSize int
	mergeSize int
	ropeType string
	size int
	replicaID string
}
/*
**
Initializes a new rope with two empty leaf nodes
*/
func NewRope(maximumChunkLength int, splitRatio float64, mergeRatio float64, ropeType string, replicaID string) *Rope {
	root := NewInnerNode(0, 0, nil, nil, nil)
	root.SetLeft(NewLeafNode([]Block{NewCRDTBlk(NewClockOffset(VectorClock{replicaID:0},0),NewRopeValue(ropeType,""),false)}, root))
	root.SetRight(NewLeafNode([]Block{NewCRDTBlk(NewClockOffset(VectorClock{replicaID:0},0),NewRopeValue(ropeType,""),false)}, root))
	splitSize := max(1, int(splitRatio*float64(maximumChunkLength)))
	mergeSize := max(1, int(mergeRatio*float64(maximumChunkLength)))

	return &Rope{
		root,
		maximumChunkLength,
		splitSize,
		mergeSize,
		ropeType,
		0,
		replicaID,
	}
}
func (r *Rope) ReplicaID()string {
	return r.replicaID
}
func (r *Rope) Root()*InnerNode {
	return r.root
}
func (r *Rope) NewRopeValue(input string) RopeValue {
	return NewRopeValue(r.ropeType, input)
}
func (r *Rope) CopyRopeValue(ropeValue RopeValue) RopeValue {
	return CopyRopeValue(r.ropeType, ropeValue)
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