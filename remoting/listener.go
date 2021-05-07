package remoting

// EventType means SourceObjectEventType
type EventType int

const (
	// EventTypeAdd means add event
	EventTypeAdd = iota
	// EventTypeDel means del event
	EventTypeDel
	// EventTypeUpdate means update event
	EventTypeUpdate
)
