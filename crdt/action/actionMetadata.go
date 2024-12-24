package action


type InsertionMetadata struct {
	content string
	index int
}
func NewInsertion(inputs ...any)any{
	if len(inputs)==2 {
		if content,ok:=inputs[0].(string);ok{
			if index,ok:=inputs[1].(int);ok{
				return &InsertionMetadata{content,index} 
			}
		}
	}
	return nil
}
type DeletionMetadata struct {
	length int
	index int
}
func NewDeletion(inputs... any)any {
	if len(inputs)==2{
		if length,ok:=inputs[0].(int);ok{
			if index,ok:=inputs[0].(int);ok{
				return &DeletionMetadata{length,index}
			}
		}
	}
	return nil
}
///////////// continue the rest of the types