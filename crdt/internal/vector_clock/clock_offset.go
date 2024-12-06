package vector_clock

type ClockOffset struct {
	offset      int
	vectorClock VectorClock
}

func NewClockOffset(vc VectorClock, offset int) *ClockOffset {
	return &ClockOffset{
		vectorClock: vc,
		offset:      offset,
	}
}
func (c *ClockOffset) Copy() *ClockOffset {
	return &ClockOffset{
		vectorClock: c.vectorClock.Copy(),
		offset:      c.offset,
	}
}

func (c *ClockOffset) Offset() int {
	return c.offset
}
func (c *ClockOffset) VectorClock() VectorClock {
	return c.vectorClock
}
func (c *ClockOffset) Compare (co *ClockOffset)int{
	return c.vectorClock.Compare(co.vectorClock)
}
func (c *ClockOffset) Merge (co *ClockOffset) *ClockOffset{
	return NewClockOffset(c.vectorClock.Merge(co.vectorClock),0);
}
