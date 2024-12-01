package value

type constructors = struct {
	constructor     func(string) BlockValue
	copyConstructor func(BlockValue) BlockValue
}

var BlockValueRegistry = make(map[string]constructors)

type BlockValue interface {
	Update(index int, input string, deletedLength int) (BlockValue, error)
	Len() int
	Split(startIndex int, endIndex int) BlockValue
	SplitFrom(startIndex int) BlockValue
	SplitTo(endIndex int) BlockValue
	String() string
}

func NewBlockValue(typename string, value string) BlockValue {
	BlockValueType, ok := BlockValueRegistry[typename]
	if ok {
		return BlockValueType.constructor(value)
	}
	return nil
}
func CopyBlockValue(typename string, value BlockValue) BlockValue {
	BlockValueType, ok := BlockValueRegistry[typename]
	if ok {
		return BlockValueType.copyConstructor(value)
	}
	return nil
}
func RegisterNewRopeType(typename string, constructor func(string) BlockValue, copyConstructor func(BlockValue) BlockValue) {
	BlockValueRegistry[typename] = constructors{constructor, copyConstructor}
}

func RegisterRopeTypes() {
	RegisterNewRopeType("ropeBuffer", NewRopeBuffer, CopyRopeBuffer)
}
func init() {
	RegisterRopeTypes()
}