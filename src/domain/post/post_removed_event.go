package domainpost

type PostRemovedEvent struct{}

func (event *PostRemovedEvent) Type() PostEventType {
	return REMOVED
}

func NewPostRemovedEvent() *PostRemovedEvent {
	return &PostRemovedEvent{}
}
