package event


type EventMetadataConstructor func(inputs ...any) EventMetadata

type EventType int

const (
	Insert EventType = iota
	Delete
)

var EventMetadataRegistry map[EventType]EventMetadataConstructor = make(map[EventType]EventMetadataConstructor)

func (et *EventType) String()string{
	switch *et{
	case Insert:
		return "Insert"
	case Delete:
		return "Delete"
	default:
		return "Undefined"
	}
}
func registryNewEventMetadata(eventType EventType, eventMetadata EventMetadataConstructor) bool {
	_, ok := EventMetadataRegistry[eventType]
	if ok || eventMetadata == nil {
		return false
	}
	EventMetadataRegistry[eventType] = eventMetadata
	return true
}
func initializeEventMetadataRegistry() {
	registryNewEventMetadata(Insert, NewInsertionEventMetadata)
	registryNewEventMetadata(Delete, NewDeletionEventMetadata)
}
func init() {
	initializeEventMetadataRegistry()
}