package domainchat

type EventType string

const (
	CREATED      EventType = "CREATED"
	MESSAGE_SENT EventType = "MESSAGE_SENT"
)

type Event interface {
	Type() EventType
}
