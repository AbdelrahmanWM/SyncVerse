package node

import (
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/block_ds"
)

type RopeNode interface {
	Right() RopeNode
	Left() RopeNode
	Parent() RopeNode
	Weight() int
	SetRight(node RopeNode)
	SetLeft(node RopeNode)
	SetParent(node RopeNode)
	SetWeight(w int)
	// String(addDeleted bool,blockSeparator string) string
}

type Node struct {
	weight int
	right  RopeNode
	left   RopeNode
	parent RopeNode
}
type InnerNode struct {
	leftWeight     int
	realLeftWeight int
	Node
}
type LeafNode struct {
	blocks BlockDS
	Node
}

func (r *Node) Right() RopeNode {
	return r.right
}
func (r *Node) Left() RopeNode {
	return r.left
}
func (r *Node) Parent() RopeNode {
	return r.parent
}
func (r *Node) Weight() int {
	return r.weight
}
func (r *Node) SetWeight(w int) {
	r.weight = w
}
func (r *Node) SetRight(node RopeNode) {
	r.right = node
}
func (r *Node) SetLeft(node RopeNode) {
	r.left = node
}
func (r *Node) SetParent(node RopeNode) {
	r.parent = node
}

//////////////////////////////////////////////////////////////

func (r *InnerNode) LeftWeight() int {
	return r.leftWeight
}

func (r *InnerNode) SetLeftWeight(w int) {
	r.leftWeight = w
}
func (r *InnerNode) RealLeftWeight() int {
	return r.realLeftWeight
}
func (r *InnerNode) SetRealLeftWeight(w int) {
	r.realLeftWeight = w
}

// //////////////////////////////////////////////////
func (r *LeafNode) Blocks() BlockDS {
	return r.blocks
}
func (r *LeafNode) SetBlocks(blocks BlockDS) {
	r.blocks = blocks
}

func (r *LeafNode) Weight() int {
	return r.blocks.Size()
}

func (r *LeafNode) String(addDeleted bool, blockSeparator string) string {
	return r.blocks.String(addDeleted, blockSeparator)
}
func NewInnerNode(leftWeight, realLeftWeight, weight int, left, right, parent RopeNode) *InnerNode {
	return &InnerNode{
		leftWeight,
		realLeftWeight,
		*NewNode(weight, left, right, parent),
	}
}
func NewLeafNode(blocks BlockDS, parent RopeNode) *LeafNode {
	return &LeafNode{
		blocks,
		*NewNode(blocks.Size(), nil, nil, parent),
	}
}

func NewNode(weight int, left, right, parent RopeNode) *Node {
	return &Node{
		weight,
		right,
		left,
		parent,
	}
}
