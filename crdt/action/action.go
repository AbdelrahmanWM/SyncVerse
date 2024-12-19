package action

type Action struct {
	Kind    ActionCode
	Metadata string
	By      string
	index   int
	content string
}

type ActionCode int

const (
	Insert ActionCode = iota
	Delete
	MainHeading
	SubHeading
	SubSubHeading
	Bold
	Italic
	BulletPoint
	Link
	Color
	BackgroundColor
)

func NewAction(kind ActionCode, metadata string, by string, index int, content string) *Action {
	return &Action{kind,metadata, by, index, content}
}
