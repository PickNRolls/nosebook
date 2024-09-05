package domainmessage

type CreatedEvent struct {
}

func (event *CreatedEvent) Type() EventType {
	return CREATED
}

func NewCreatedEvent() *CreatedEvent {
	return &CreatedEvent{}
}
