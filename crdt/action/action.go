package action

import "github.com/AbdelrahmanWM/SyncVerse/crdt/global"

type Action struct {
	Kind      ActionCode
	UserID    global.UserID
	ReplicaID global.ReplicaID
	Metadata  any
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

type ActionConstructor func(inputs ...any) any

var ActionMetadataRegistry map[ActionCode]ActionConstructor = make(map[ActionCode]ActionConstructor)

func NewAction(kind ActionCode, userID global.UserID, replicaID global.ReplicaID, metadata any) *Action {
	return &Action{kind, userID, replicaID, metadata}
}
func registerMetadata(actionCode ActionCode, metadata ActionConstructor) bool {
	_, ok := ActionMetadataRegistry[actionCode]
	if ok || metadata == nil {
		return false
	}
	ActionMetadataRegistry[actionCode] = metadata
	return true
}
func initializeMetadataMap() {
	registerMetadata(Insert, NewInsertion)
	registerMetadata(Delete, NewDeletion)
}
func init() {
	initializeMetadataMap()
}
