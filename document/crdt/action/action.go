package action

import (
	"strings"
)

type Action struct {
	Kind      ActionCode
	Metadata  ActionMetadata
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

func (a ActionCode) String() string {
	switch a {
	case Insert:
		return "Insert"
	case Delete:
		return "Delete"
	case MainHeading:
		return "MainHeading"
	case SubHeading:
		return "SubHeading"
	case SubSubHeading:
		return "SubSubHeading"
	case Bold:
		return "Bold"
	case Italic:
		return "Italic"
	case BulletPoint:
		return "BulletPoint"
	case Link:
		return "Link"
	case Color:
		return "Color"
	case BackgroundColor:
		return "BackgroundColor"
	default:
		return "Unknown ActionCode"
	}
}

type ActionConstructor func(inputs ...any) any

var ActionMetadataRegistry map[ActionCode]ActionConstructor = make(map[ActionCode]ActionConstructor)

func NewAction(kind ActionCode,  metadata ActionMetadata) *Action {
	return &Action{kind,  metadata}
}
func registerMetadata(actionCode ActionCode, metadata ActionConstructor) bool {
	_, ok := ActionMetadataRegistry[actionCode]
	if ok || metadata == nil {
		return false
	}
	ActionMetadataRegistry[actionCode] = metadata
	return true
}

func (a *Action) String() string {
	var results strings.Builder
	results.WriteString("Action\n")
	results.WriteString(a.Kind.String())
	results.WriteString("\n")
	results.WriteString(a.Metadata.String())
	results.WriteString("\n")
	return results.String()
}
func initializeMetadataMap() {
	registerMetadata(Insert, NewInsertion)
	registerMetadata(Delete, NewDeletion)
}
func init() {
	initializeMetadataMap()
}
