package value

import (
	"fmt"
)

type RopeBuffer []byte

func (rp *RopeBuffer) Update(index int, newContent string, deletedLength int) (BlockValue, error) {
	input := NewBlockValue(ByteBuffer, newContent)
	if index < 0 {
		return NewRopeBuffer(""), fmt.Errorf("[ERROR] Invalid Index %d", index)
	} else if index >= rp.Len() {
		err := rp.append(input)
		return rp, err
	}

	if buf, ok := input.(*RopeBuffer); ok {
		endOfDeletion := min(index+deletedLength, rp.Len())

		*rp = append((*rp)[:index], append(*buf, (*rp)[endOfDeletion:]...)...)
		return rp, nil
	}
	return NewRopeBuffer(""), fmt.Errorf("[ERROR] Invalid Input Type")
}
func (rp *RopeBuffer) append(input BlockValue) error {
	if buf, ok := input.(*RopeBuffer); ok {
		(*rp) = append(*rp, (*buf)...)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Invalid Input Type")
	}
}
func (rp *RopeBuffer) Split(startIndex int, endIndex int) BlockValue {
	value := (*rp)[startIndex:endIndex]
	return &value
}

func (rp *RopeBuffer) SplitFrom(startIndex int) BlockValue {
	value := (*rp)[startIndex:]
	return &value
}
func (rp *RopeBuffer) SplitTo(endIndex int) BlockValue {
	value := (*rp)[:endIndex]
	return &value
}
func (rp *RopeBuffer) Len() int {
	return len([]byte(*rp))
}

func NewRopeBuffer(b string) BlockValue {
	rb := RopeBuffer(b)
	return &rb
}
func CopyRopeBuffer(rp BlockValue) BlockValue {
	rb := RopeBuffer(rp.String())
	return &rb
}
func (rp *RopeBuffer) String() string {
	return string(*rp)
}
