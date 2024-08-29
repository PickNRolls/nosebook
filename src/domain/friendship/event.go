package domainfriendship

type EventType string

const (
	CREATED  EventType = "CREATED"
	ACCEPTED EventType = "ACCEPTED"
	DENIED   EventType = "DENIED"
	REMOVED  EventType = "REMOVED"
	VIEWED   EventType = "VIEWED"
)

type Event interface {
	Type() EventType
}
