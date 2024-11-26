package block

import (
	"errors"
	. "github.com/AbdelrahmanWM/SyncVerse/crdt/internal/rope/value"
)

type CRDTBlk struct {
	clockOffset *ClockOffset
	content     RopeValue
	deleted     bool
}

func NewCRDTBlk(clockOffset *ClockOffset, content RopeValue, deleted bool) *CRDTBlk {
	return &CRDTBlk{clockOffset: clockOffset, content: content, deleted: deleted}
}

func (c *CRDTBlk) Len() int {
	return c.content.Len()
}
func (c *CRDTBlk) Split(index int) (Block, Block) {
	leftContent, rightContent := c.content.SplitTo(index), c.content.SplitFrom(index)

	return NewCRDTBlk(c.clockOffset.Copy(), leftContent, c.deleted), NewCRDTBlk(NewClockOffset(c.clockOffset.vectorClock.Copy(), index), rightContent, c.deleted)
}

func (c *CRDTBlk)String() string{
	return c.content.String()
}
func (c *CRDTBlk)IsDeleted()bool{
	return c.deleted
}
func (c *CRDTBlk) Compare(b Block) (int, error) {
	crdtBlk, ok := b.(*CRDTBlk)
	if !ok {
		return 0, errors.New("[ERROR] Incomparable block types")
	}
	return c.clockOffset.Compare(crdtBlk.clockOffset), nil
}

func (c*CRDTBlk)Offset() int{
	return c.clockOffset.Offset()
}
func (c *CRDTBlk) Delete() {
	c.deleted = true
}
