package block

type Block interface {
	Len() int
	Split(index int) (Block, Block)
	String() string
	IsDeleted()bool
	Compare(b Block) (int, error)
	Offset()int
	Delete()
}