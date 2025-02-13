package block

import (
	"strconv"
	"strings"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/action"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/format"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/rope/value"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
)

type Block struct {
	clockOffset *ClockOffset
	content     BlockValue
	blockType   ValueType
	deleted     bool
	formats     map[ActionCode]Format
}

func NewBlock(clockOffset *ClockOffset, content string, blockValueType ValueType, deleted bool) *Block {
	return &Block{clockOffset: clockOffset, content: NewBlockValue(blockValueType, content), blockType: blockValueType, deleted: deleted, formats: make(map[ActionCode]Format, 1)}
}
func CopyBlock(block *Block) *Block {
	return &Block{clockOffset: block.clockOffset.Copy(), content: CopyBlockValue(block.blockType, block.content), blockType: block.blockType, deleted: block.deleted}
}
func (b *Block) NewBlockValue(input string) BlockValue {
	return NewBlockValue(b.blockType, input)
}
func (b *Block) CopyRopeValue(ropeValue BlockValue) BlockValue {
	return CopyBlockValue(b.blockType, ropeValue)
}
func (c *Block) Len() int {
	return c.content.Len()
}
func (c *Block) Split(index int) (*Block, *Block) {
	if index == 0 {
		return nil, c
	}
	leftContent, rightContent := c.content.SplitTo(index), c.content.SplitFrom(index)
	return NewBlock(c.clockOffset.Copy(), leftContent.String(), c.blockType, c.deleted), NewBlock(NewClockOffset(c.clockOffset.VectorClock().Copy(), index), rightContent.String(), c.blockType, c.deleted)
}

func (c *Block) Content() string {
	return c.content.String()
}
func (c *Block) String() string {
	var results strings.Builder
	results.WriteString("Block\n")
	results.WriteString(c.Content())
	results.WriteString(" ")
	results.WriteString(c.clockOffset.String())
	results.WriteString(" ")
	results.WriteString(c.blockType.String())
	results.WriteString(" ")
	results.WriteString(strconv.FormatBool(c.deleted))
	return results.String()
}
func (c *Block) IsDeleted() bool {
	return c.deleted
}
func (c *Block) Compare(b *Block) (int, error) {
	return c.clockOffset.Compare(b.clockOffset), nil
}
func (c *Block) ClockOffset() *ClockOffset {
	return c.clockOffset
}

func (c *Block) Offset() int {
	return c.clockOffset.Offset()
}
func (c *Block) Delete() {
	c.deleted = true
}
func (c *Block) ContainsOffset(offset int) bool {
	return offset >= c.Offset() && offset < c.Offset()+c.content.Len()
}
func (c *Block) HasVectorClock(vectorClock VectorClock) bool {
	return c.clockOffset.VectorClock().Equals(vectorClock)
}
func (c *Block) CompareHashes(c2 *Block) int {
	return c.clockOffset.CompareHashes(c2.clockOffset)
}
func (c *Block) AddFormatting(format Format) {
	c.formats[format.Kind] = format
}
func (c *Block) RemoveFormatting(format Format) {
	delete(c.formats, format.Kind)
}
func (c *Block) FormatExists(actionCode ActionCode) bool {
	_, ok := c.formats[actionCode]
	return ok
}
