package action

import (
	"strconv"
	"strings"
)

type ActionMetadata interface {
	String() string
}
type InsertionMetadata struct {
	Content string
	Index   int
}

func NewInsertion(inputs ...any) any {
	if len(inputs) == 2 {
		if content, ok := inputs[0].(string); ok {
			if index, ok := inputs[1].(int); ok {
				return &InsertionMetadata{content, index}
			}
		}
	}
	return nil
}
func (iam *InsertionMetadata) String() string {
	var result strings.Builder
	result.WriteString(iam.Content)
	result.WriteString(" ")
	result.WriteString("Index: ")
	result.WriteString(strconv.Itoa(iam.Index))
	return result.String()
}

type DeletionMetadata struct {
	Length int
	Index  int
}

func NewDeletion(inputs ...any) any {
	if len(inputs) == 2 {
		if length, ok := inputs[0].(int); ok {
			if index, ok := inputs[0].(int); ok {
				return &DeletionMetadata{length, index}
			}
		}
	}
	return nil
}
func (iam *DeletionMetadata) String() string {
	var result strings.Builder
	result.WriteString("Deletion length: ")
	result.WriteString(strconv.Itoa(iam.Length))
	result.WriteString(" ")
	result.WriteString("Index: ")
	result.WriteString(strconv.Itoa(iam.Index))
	return result.String()
}
///////////// continue the rest of the types