package domainmessage

type RemovedEvent struct{}

func (event *RemovedEvent) Type() EventType {
	return REMOVED
}

func NewRemovedEvent() *RemovedEvent {
	return &RemovedEvent{}
}
