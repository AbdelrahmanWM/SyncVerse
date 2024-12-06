package action

type Action struct {
	Kind    ActionCode
	By      string
	index   int
	content string
}

type ActionCode int

const (
	Insert ActionCode = iota
	Delete
	Update
)

func NewAction(kind ActionCode, by string, index int, content string) *Action {
	return &Action{kind, by, index, content}
}
