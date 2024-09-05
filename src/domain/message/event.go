package domainmessage

type EventType string

const (
	CREATED EventType = "CREATED"
	REMOVED EventType = "REMOVED"
)

type Event interface {
	Type() EventType
}
