package value

type constructors = struct {
	constructor     func(string) RopeValue
	copyConstructor func(RopeValue) RopeValue
}

var ropeValueRegistry = make(map[string]constructors)

type RopeValue interface {
	Update(index int, input string, deletedLength int) (RopeValue, error)
	Len() int
	Split(startIndex int, endIndex int) RopeValue
	SplitFrom(startIndex int) RopeValue
	SplitTo(endIndex int) RopeValue
	String() string
}

func NewRopeValue(typename string, value string) RopeValue {
	ropeValueType, ok := ropeValueRegistry[typename]
	if ok {
		return ropeValueType.constructor(value)
	}
	return nil
}
func CopyRopeValue(typename string, value RopeValue) RopeValue {
	ropeValueType, ok := ropeValueRegistry[typename]
	if ok {
		return ropeValueType.copyConstructor(value)
	}
	return nil
}
func RegisterNewRopeType(typename string, constructor func(string) RopeValue, copyConstructor func(RopeValue) RopeValue) {
	ropeValueRegistry[typename] = constructors{constructor, copyConstructor}
}

func RegisterRopeTypes() {
	RegisterNewRopeType("ropeBuffer", NewRopeBuffer, CopyRopeBuffer)
}
func init() {
	RegisterRopeTypes()
}