package posts

type PostCreatedEvent struct {
}

func (event *PostCreatedEvent) Type() PostEventType {
	return CREATED
}

func NewPostCreatedEvent() *PostCreatedEvent {
	return &PostCreatedEvent{}
}
