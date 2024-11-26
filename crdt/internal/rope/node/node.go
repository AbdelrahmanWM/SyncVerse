package node

import (
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/block"
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
}

type Node struct {
	weight int
	right  RopeNode
	left   RopeNode
	parent RopeNode
}
type InnerNode struct {
	leftWeight int
	Node
}
type LeafNode struct {
	blocks []Block
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

// //////////////////////////////////////////////////
func (r *LeafNode) Blocks() []Block {
	return r.blocks
}
func (r *LeafNode) SetBlocks(blocks []Block) {
	r.blocks = blocks
}

func (r *LeafNode) Weight() int {
	return len(r.blocks)
}

func NewInnerNode(leftWeight, weight int, left, right, parent RopeNode) *InnerNode {
	return &InnerNode{
		leftWeight,
		*NewNode(weight, left, right, parent),
	}
}
func NewLeafNode(blocks []Block, parent RopeNode) *LeafNode {
	return &LeafNode{
		blocks,
		*NewNode(len(blocks), nil, nil, parent),
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
