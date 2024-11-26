package value

import (
	"fmt"
)

type RopeBuffer []byte

func (rp *RopeBuffer) Update(index int, newContent string, deletedLength int) (RopeValue, error) {
	input := NewRopeValue("ropeBuffer", newContent)
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
func (rp *RopeBuffer) append(input RopeValue) error {
	if buf, ok := input.(*RopeBuffer); ok {
		(*rp) = append(*rp, (*buf)...)
		return nil
	} else {
		return fmt.Errorf("[ERROR] Invalid Input Type")
	}
}
func (rp *RopeBuffer) Split(startIndex int, endIndex int) RopeValue {
	value := (*rp)[startIndex:endIndex]
	return &value
}

func (rp *RopeBuffer) SplitFrom(startIndex int) RopeValue {
	value := (*rp)[startIndex:]
	return &value
}
func (rp *RopeBuffer) SplitTo(endIndex int) RopeValue {
	value := (*rp)[:endIndex]
	return &value
}
func (rp *RopeBuffer) Len() int {
	return len([]byte(*rp))
}

func NewRopeBuffer(b string) RopeValue {
	rb := RopeBuffer(b)
	return &rb
}
func CopyRopeBuffer(rp RopeValue) RopeValue {
	rb := RopeBuffer(rp.String())
	return &rb
}
func (rp *RopeBuffer) String() string {
	return string(*rp)
}
